package dbeval

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

type PGXStdlib struct {
	db *sql.DB
}

func (p *PGXStdlib) Connect(ds string, connLifetime time.Duration, idleConns, openConns int) {
	if p.db != nil {
		check(p.db.Close())
		p.db = nil
	}
	var err error
	p.db, err = sql.Open("pgx", ds)
	check(err)
	p.db.SetConnMaxLifetime(connLifetime)
	p.db.SetMaxIdleConns(idleConns)
	p.db.SetMaxOpenConns(openConns)
}

func (p *PGXStdlib) CreateDatabase() {
	_, err := p.db.Exec(createdb)
	check(err)
}

func (p *PGXStdlib) DropDatabase() {
	_, err := p.db.Exec(dropdb)
	check(err)
}

func (p *PGXStdlib) CreateSchema() {
	_, err := p.db.Exec(schema)
	check(err)
}

// copyFrom borrows the underlying *pgx.Conn from the database/sql pool, the documented way to reach pgx-specific
// features such as COPY through database/sql.
func (p *PGXStdlib) copyFrom(table string, columns []string, src pgx.CopyFromSource) {
	ctx := context.Background()
	conn, err := p.db.Conn(ctx)
	check(err)
	defer conn.Close()
	check(conn.Raw(func(driverConn any) error {
		_, err := driverConn.(*stdlib.Conn).Conn().CopyFrom(ctx, pgx.Identifier{table}, columns, src)
		return err
	}))
}

func (p *PGXStdlib) InsertAuthors(as []*Author) {
	p.copyFrom("authors", []string{"id", "name"}, pgx.CopyFromSlice(len(as), func(i int) ([]any, error) {
		return []any{as[i].ID, as[i].Name}, nil
	}))
}

func (p *PGXStdlib) InsertArticles(as []*Article) {
	p.copyFrom("articles", []string{"id", "title", "body", "published_at"}, pgx.CopyFromSlice(len(as), func(i int) ([]any, error) {
		return []any{as[i].ID, as[i].Title, as[i].Body, as[i].PublishedAt}, nil
	}))
}

func (p *PGXStdlib) FindAuthorByID(id int64) *Author {
	a := &Author{}
	check(p.db.QueryRow(`SELECT id, name FROM authors WHERE id = $1`, id).Scan(&a.ID, &a.Name))
	return a
}

func (p *PGXStdlib) FindAuthorsByName(name string) []*Author {
	var as []*Author
	rows, err := p.db.Query(`SELECT id, name FROM authors WHERE name = $1`, name)
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

func (p *PGXStdlib) RecentArticles(n int) []*Article {
	var as []*Article
	rows, err := p.db.Query(`SELECT id, title, body, published_at FROM articles ORDER BY published_at DESC LIMIT $1`, n)
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

var _ Implementation = &PGXStdlib{}
