package postings

import (
	"context"

	"github.com/kimseogyu/portfolio/backend/internal/db"
	"gorm.io/gorm"
)

type Repository interface {
	Save(ctx context.Context, data *Posting) error
	Delete(ctx context.Context, id int) error
	FindAll(ctx context.Context, pagination *db.CursorBasedPagination) (*db.CursorBasedPaginationResponse[Posting], error)
	FindOneByID(ctx context.Context, id int) (*Posting, error)
	Update(ctx context.Context, id int, data Posting) error
	IncrementViewCount(ctx context.Context, id int) error
}

type postingsGormRepository struct {
	db *gorm.DB
}

var _ Repository = &postingsGormRepository{}

// Save implements db.Repository.
func (p *postingsGormRepository) Save(ctx context.Context, data *Posting) error {
	tx := p.db.WithContext(ctx)
	if err := tx.Save(data).Error; err != nil {
		return err
	}

	return nil
}

// Delete implements db.Repository.
func (p *postingsGormRepository) Delete(ctx context.Context, id int) error {
	tx := p.db.WithContext(ctx)
	posting := &Posting{}
	if err := tx.First(posting, id).Error; err != nil {
		return err
	}

	return tx.Model(posting).Update("deleted_at", gorm.Expr("NOW()")).Error
}

// FindAll implements db.Repository.
func (p *postingsGormRepository) FindAll(ctx context.Context, pagination *db.CursorBasedPagination) (*db.CursorBasedPaginationResponse[Posting], error) {
	tx := p.db.WithContext(ctx)

	var total int64
	err := tx.Model(&Posting{}).Count(&total).Error
	if err != nil {
		return nil, err
	}

	tx = tx.Model(&Posting{})

	if pagination.Cursor != nil {
		tx = tx.Where("id > ?", *pagination.Cursor)
	}

	if pagination.Limit != nil {
		tx = tx.Limit(int(*pagination.Limit))
	}

	var postings []Posting
	err = tx.Find(&postings).Error
	if err != nil {
		return nil, err
	}

	nextCursor := int64(0)
	if len(postings) > 0 {
		nextCursor = int64(postings[len(postings)-1].ID)
	}

	return &db.CursorBasedPaginationResponse[Posting]{
		Data:       postings,
		NextCursor: &nextCursor,
		HasNext:    len(postings) == int(*pagination.Limit),
		Total:      int(total),
	}, nil
}

// FindOneByID implements db.Repository.
func (p *postingsGormRepository) FindOneByID(ctx context.Context, id int) (*Posting, error) {
	tx := p.db.WithContext(ctx)

	var posting Posting
	err := tx.Where("id = ?", id).First(&posting).Error
	if err != nil {
		return nil, err
	}

	return &posting, nil
}

// Update implements db.Repository.
func (p *postingsGormRepository) Update(ctx context.Context, id int, data Posting) error {
	tx := p.db.WithContext(ctx)
	return tx.Model(&Posting{}).Where("id = ?", id).Updates(data).Error
}

// IncrementViewCount implements db.Repository.
func (p *postingsGormRepository) IncrementViewCount(ctx context.Context, id int) error {
	tx := p.db.WithContext(ctx)
	return tx.Model(&Posting{}).Where("id = ?", id).
		Update("view_count", gorm.Expr("view_count + ?", 1)).Error
}

func NewRepository(db *gorm.DB) *postingsGormRepository {
	return &postingsGormRepository{
		db: db,
	}
}
