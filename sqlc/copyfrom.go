// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: copyfrom.go

package sqlc

import (
	"context"
)

// iteratorForInsertArticles implements pgx.CopyFromSource.
type iteratorForInsertArticles struct {
	rows                 []InsertArticlesParams
	skippedFirstNextCall bool
}

func (r *iteratorForInsertArticles) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForInsertArticles) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].ID,
		r.rows[0].Title,
		r.rows[0].Body,
		r.rows[0].PublishedAt,
	}, nil
}

func (r iteratorForInsertArticles) Err() error {
	return nil
}

func (q *Queries) InsertArticles(ctx context.Context, arg []InsertArticlesParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"articles"}, []string{"id", "title", "body", "published_at"}, &iteratorForInsertArticles{rows: arg})
}

// iteratorForInsertAuthors implements pgx.CopyFromSource.
type iteratorForInsertAuthors struct {
	rows                 []InsertAuthorsParams
	skippedFirstNextCall bool
}

func (r *iteratorForInsertAuthors) Next() bool {
	if len(r.rows) == 0 {
		return false
	}
	if !r.skippedFirstNextCall {
		r.skippedFirstNextCall = true
		return true
	}
	r.rows = r.rows[1:]
	return len(r.rows) > 0
}

func (r iteratorForInsertAuthors) Values() ([]interface{}, error) {
	return []interface{}{
		r.rows[0].ID,
		r.rows[0].Name,
	}, nil
}

func (r iteratorForInsertAuthors) Err() error {
	return nil
}

func (q *Queries) InsertAuthors(ctx context.Context, arg []InsertAuthorsParams) (int64, error) {
	return q.db.CopyFrom(ctx, []string{"authors"}, []string{"id", "name"}, &iteratorForInsertAuthors{rows: arg})
}
