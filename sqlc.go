package dbeval

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/davars/dbeval/sqlc"
)

type SQLC struct {
	db *pgxpool.Pool
	q  *sqlc.Queries
}

func (p *SQLC) Connect(ds string, connLifetime time.Duration, idleConns, openConns int) {
	if p.db != nil {
		p.db.Close()
		p.db = nil
	}
	cc, err := pgxpool.ParseConfig(ds)
	check(err)
	if connLifetime > 0 {
		cc.MaxConnLifetime = connLifetime
	}
	cc.MinConns = int32(idleConns)
	cc.MaxConns = int32(openConns)
	p.db, err = pgxpool.NewWithConfig(context.Background(), cc)
	p.q = sqlc.New(p.db)
}

func (p *SQLC) CreateDatabase() {
	_, err := p.db.Exec(context.Background(), createdb)
	check(err)
}

func (p *SQLC) DropDatabase() {
	_, err := p.db.Exec(context.Background(), dropdb)
	check(err)
}

func (p *SQLC) CreateSchema() {
	_, err := p.db.Exec(context.Background(), schema)
	check(err)
}

func (p *SQLC) InsertAuthors(as []*Author) {
	insert := make([]sqlc.InsertAuthorsParams, len(as))
	for i, a := range as {
		insert[i].ID = a.ID
		insert[i].Name = a.Name
	}

	ctx := context.Background()
	tx, err := p.db.Begin(ctx)
	defer func() {
		if err == nil {
			check(tx.Commit(ctx))
		} else {
			check(tx.Rollback(ctx))
		}
	}()
	check(err)

	_, err = p.q.WithTx(tx).InsertAuthors(context.Background(), insert)
	check(err)
}

func (p *SQLC) InsertArticles(as []*Article) {
	insert := make([]sqlc.InsertArticlesParams, len(as))
	for i, a := range as {
		insert[i].ID = a.ID
		insert[i].Body = a.Body
		insert[i].Title = a.Title
		check((&insert[i].PublishedAt).Scan(a.PublishedAt))
	}

	ctx := context.Background()
	tx, err := p.db.Begin(ctx)
	defer func() {
		if err == nil {
			check(tx.Commit(ctx))
		} else {
			check(tx.Rollback(ctx))
		}
	}()
	check(err)

	_, err = p.q.WithTx(tx).InsertArticles(context.Background(), insert)
	check(err)
}

func (p *SQLC) FindAuthorByID(id int64) *Author {
	a, err := p.q.FindAuthorByID(context.Background(), id)
	check(err)
	return &Author{
		ID:   a.ID,
		Name: a.Name,
	}
}

func (p *SQLC) FindAuthorsByName(name string) []*Author {
	as, err := p.q.FindAuthorsByName(context.Background(), name)
	check(err)
	ret := make([]*Author, len(as))
	for i, a := range as {
		ret[i] = &Author{
			ID:   a.ID,
			Name: a.Name,
		}
	}
	return ret
}

func (p *SQLC) RecentArticles(n int) []*Article {
	as, err := p.q.RecentArticles(context.Background(), int32(n))
	check(err)
	ret := make([]*Article, len(as))
	for i, a := range as {
		if !a.PublishedAt.Valid {
			check(fmt.Errorf("invalid PublishedAt"))
		}
		ret[i] = &Article{
			ID:          a.ID,
			Body:        a.Body,
			Title:       a.Title,
			PublishedAt: a.PublishedAt.Time,
		}
	}
	return ret
}

var _ Implementation = &SQLC{}
