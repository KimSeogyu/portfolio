CREATE TABLE IF NOT EXISTS postings (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NULL,
    author_id VARCHAR(255) NOT NULL,
    author_name VARCHAR(255) NOT NULL,
    comment_count INT NOT NULL DEFAULT 0,
    view_count INT NOT NULL DEFAULT 0,
    tags VARCHAR(255)[] NULL,
    status SMALLINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_postings_author_id ON postings (author_id);
CREATE INDEX idx_postings_author_name ON postings (author_name);
CREATE INDEX idx_postings_tags ON postings USING GIN (tags);
CREATE INDEX idx_postings_deleted_at ON postings (deleted_at);
CREATE INDEX idx_postings_status ON postings (status);
CREATE INDEX idx_postings_created_at ON postings (created_at);

COMMENT ON TABLE postings IS '게시물 테이블';
COMMENT ON COLUMN postings.id IS '게시물 ID';
COMMENT ON COLUMN postings.title IS '게시물 제목';
COMMENT ON COLUMN postings.content IS '게시물 내용';
COMMENT ON COLUMN postings.author_id IS '게시물 작성자 ID';
COMMENT ON COLUMN postings.author_name IS '게시물 작성자 이름';
COMMENT ON COLUMN postings.comment_count IS '게시물 댓글 수';
COMMENT ON COLUMN postings.view_count IS '게시물 조회 수';
COMMENT ON COLUMN postings.tags IS '게시물 태그';
COMMENT ON COLUMN postings.status IS '게시물 상태 (0: 비활성, 1: 활성, 2: Draft, 3: 삭제)';
COMMENT ON COLUMN postings.created_at IS '게시물 생성 시간';
COMMENT ON COLUMN postings.updated_at IS '게시물 수정 시간';
COMMENT ON COLUMN postings.deleted_at IS '게시물 삭제 시간';

CREATE TABLE IF NOT EXISTS comments (
    id BIGSERIAL PRIMARY KEY,
    post_id BIGINT NOT NULL,
    parent_id BIGINT,
    author_id VARCHAR(255) NOT NULL,
    author_name VARCHAR(255) NOT NULL,
    children_count INT NOT NULL DEFAULT 0,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    status SMALLINT NOT NULL DEFAULT 0
);

CREATE INDEX idx_comments_post_id ON comments (post_id);
CREATE INDEX idx_comments_parent_id ON comments (parent_id);
CREATE INDEX idx_comments_author_id ON comments (author_id);
CREATE INDEX idx_comments_author_name ON comments (author_name);
CREATE INDEX idx_comments_status ON comments (status);
CREATE INDEX idx_comments_created_at ON comments (created_at);
CREATE INDEX idx_comments_updated_at ON comments (updated_at);
CREATE INDEX idx_comments_deleted_at ON comments (deleted_at);

COMMENT ON TABLE comments IS '댓글 테이블';
COMMENT ON COLUMN comments.id IS '댓글 ID';
COMMENT ON COLUMN comments.post_id IS '게시물 ID';
COMMENT ON COLUMN comments.parent_id IS '부모 댓글 ID';
COMMENT ON COLUMN comments.author_id IS '댓글 작성자 ID';
COMMENT ON COLUMN comments.author_name IS '댓글 작성자 이름';
COMMENT ON COLUMN comments.children_count IS '댓글 자식 수';
COMMENT ON COLUMN comments.content IS '댓글 내용';
COMMENT ON COLUMN comments.created_at IS '댓글 생성 시간';
COMMENT ON COLUMN comments.updated_at IS '댓글 수정 시간';
COMMENT ON COLUMN comments.deleted_at IS '댓글 삭제 시간';
COMMENT ON COLUMN comments.status IS '댓글 상태 (0: 비활성, 1: 활성, 2: 삭제)';
