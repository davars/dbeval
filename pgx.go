package dbeval

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Notes:
// - bulk inserts use CopyFrom; a single COPY is atomic on its own, so no explicit transaction is needed
// - pgx caches prepared statements per connection by default (QueryExecModeCacheStatement)

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
	_, err := p.db.CopyFrom(context.Background(), pgx.Identifier{"authors"}, []string{"id", "name"},
		pgx.CopyFromSlice(len(as), func(i int) ([]any, error) {
			return []any{as[i].ID, as[i].Name}, nil
		}))
	check(err)
}

func (p *PGX) InsertArticles(as []*Article) {
	_, err := p.db.CopyFrom(context.Background(), pgx.Identifier{"articles"}, []string{"id", "title", "body", "published_at"},
		pgx.CopyFromSlice(len(as), func(i int) ([]any, error) {
			return []any{as[i].ID, as[i].Title, as[i].Body, as[i].PublishedAt}, nil
		}))
	check(err)
}

func (p *PGX) FindAuthorByID(id int64) *Author {
	a := &Author{}
	check(p.db.QueryRow(context.Background(), `SELECT id, name FROM authors WHERE id = $1`, id).Scan(&a.ID, &a.Name))
	return a
}

func (p *PGX) FindAuthorsByName(name string) []*Author {
	rows, _ := p.db.Query(context.Background(), `SELECT id, name FROM authors WHERE name = $1`, name)
	as, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*Author, error) {
		a := &Author{}
		return a, row.Scan(&a.ID, &a.Name)
	})
	check(err)
	return as
}

func (p *PGX) RecentArticles(n int) []*Article {
	rows, _ := p.db.Query(context.Background(), `SELECT id, title, body, published_at FROM articles ORDER BY published_at DESC LIMIT $1`, n)
	as, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*Article, error) {
		a := &Article{}
		return a, row.Scan(&a.ID, &a.Title, &a.Body, &a.PublishedAt)
	})
	check(err)
	return as
}

var _ Implementation = &PGX{}
