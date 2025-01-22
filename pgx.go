package dbeval

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PGX struct {
	db *pgxpool.Pool
}

func (p *PGX) Connect(ds string, connLifetime time.Duration, idleConns, openConns int) {
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
	check(err)
}

func (p *PGX) CreateDatabase() {
	_, err := p.db.Exec(context.Background(), createdb)
	check(err)
}

func (p *PGX) DropDatabase() {
	_, err := p.db.Exec(context.Background(), dropdb)
	check(err)
}

func (p *PGX) CreateSchema() {
	_, err := p.db.Exec(context.Background(), schema)
	check(err)
}

func (p *PGX) InsertAuthors(as []*Author) {
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
	for _, a := range as {
		_, err = tx.Exec(ctx, `INSERT INTO authors (id, name) VALUES ($1, $2)`, a.ID, a.Name)
		check(err)
	}
}

func (p *PGX) InsertArticles(as []*Article) {
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
	for _, a := range as {
		_, err = tx.Exec(ctx, `INSERT INTO articles (id, title, body, published_at) VALUES ($1, $2, $3, $4)`, a.ID, a.Title, a.Body, a.PublishedAt)
		check(err)
	}
}

func (p *PGX) FindAuthorByID(id int64) *Author {
	a := &Author{}
	check(p.db.QueryRow(context.Background(), `SELECT id, name FROM authors WHERE id = $1`, id).Scan(&a.ID, &a.Name))
	return a
}

func (p *PGX) FindAuthorsByName(name string) []*Author {
	var as []*Author
	rows, err := p.db.Query(context.Background(), `SELECT id, name FROM authors WHERE name = $1`, name)
	check(err)
	defer rows.Close()
	for rows.Next() {
		a := &Author{}
		err = rows.Scan(&a.ID, &a.Name)
		check(err)
		as = append(as, a)
	}
	check(rows.Err())
	return as
}

func (p *PGX) RecentArticles(n int) []*Article {
	var as []*Article
	rows, err := p.db.Query(context.Background(), `SELECT id, title, body, published_at FROM articles ORDER BY published_at DESC LIMIT $1`, n)
	check(err)
	defer rows.Close()
	for rows.Next() {
		a := &Article{}
		err = rows.Scan(&a.ID, &a.Title, &a.Body, &a.PublishedAt)
		check(err)
		as = append(as, a)
	}
	check(rows.Err())
	return as
}

var _ Implementation = &PGX{}
