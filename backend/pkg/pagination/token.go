package pagination

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type Token struct {
	Cursor *int64
	Limit  *int
}

const DefaultLimit = 10

func NewToken(cursor int64, limit int) *Token {
	return &Token{
		Cursor: &cursor,
		Limit:  &limit,
	}
}

func (t *Token) Encode() string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d:%d", t.Cursor, t.Limit)))
}

func FromEncodedString(encoded string) (*Token, error) {
	defaultLimit := DefaultLimit
	if encoded == "" {
		return &Token{
			Cursor: nil,
			Limit:  &defaultLimit,
		}, nil
	}

	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	var token Token
	if err := json.Unmarshal(decoded, &token); err != nil {
		return nil, err
	}

	return &token, nil
}
