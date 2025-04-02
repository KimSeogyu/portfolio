package cstore

import "fmt"

var CommentsListCacheKeyFunc = func(postID int64) string {
	return fmt.Sprintf("postings:%d:comments", postID)
}

var PostingCacheKeyFunc = func(id int64) string {
	return fmt.Sprintf("postings:%d", id)
}

var PostingListCacheKeyFunc = func(cursor int64, limit int) string {
	return fmt.Sprintf("postings:list:%d:%d", cursor, limit)
}
