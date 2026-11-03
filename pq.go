package dbeval

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

// Notes:
// - binary_parameters=yes lets pq send Parse/Bind/Execute in one round trip instead of preparing the statement
//   first, roughly halving latency for parameterized queries
// - bulk inserts use COPY, which pq exposes through a prepared statement inside a transaction

type PQ struct {
	db *sql.DB
}

func (p *PQ) Connect(ds string, connLifetime time.Duration, idleConns, openConns int) {
	if p.db != nil {
		check(p.db.Close())
		p.db = nil
	}
	var err error
	p.db, err = sql.Open("postgres", ds+" binary_parameters=yes")
	check(err)
	p.db.SetConnMaxLifetime(connLifetime)
	p.db.SetMaxIdleConns(idleConns)
	p.db.SetMaxOpenConns(openConns)
}

func (p *PQ) CreateDatabase() {
	_, err := p.db.Exec(createdb)
	check(err)
}

func (p *PQ) DropDatabase() {
	_, err := p.db.Exec(dropdb)
	check(err)
}

func (p *PQ) CreateSchema() {
	_, err := p.db.Exec(schema)
	check(err)
}

// copyIn streams rows into the table with COPY FROM STDIN.
func (p *PQ) copyIn(query string, n int, row func(i int) []any) {
	tx, err := p.db.Begin()
	check(err)
	defer tx.Rollback()
	stmt, err := tx.Prepare(query)
	check(err)
	for i := 0; i < n; i++ {
		_, err = stmt.Exec(row(i)...)
		check(err)
	}
	_, err = stmt.Exec() // flush
	check(err)
	check(stmt.Close())
	check(tx.Commit())
}

func (p *PQ) InsertAuthors(as []*Author) {
	p.copyIn(`COPY authors (id, name) FROM STDIN`, len(as), func(i int) []any {
		return []any{as[i].ID, as[i].Name}
	})
}

func (p *PQ) InsertArticles(as []*Article) {
	p.copyIn(`COPY articles (id, title, body, published_at) FROM STDIN`, len(as), func(i int) []any {
		return []any{as[i].ID, as[i].Title, as[i].Body, as[i].PublishedAt}
	})
}

func (p *PQ) FindAuthorByID(id int64) *Author {
	a := &Author{}
	check(p.db.QueryRow(`SELECT id, name FROM authors WHERE id = $1`, id).Scan(&a.ID, &a.Name))
	return a
}

func (p *PQ) FindAuthorsByName(name string) []*Author {
	var as []*Author
	rows, err := p.db.Query(`SELECT id, name FROM authors WHERE name = $1`, name)
	check(err)
	defer rows.Close()
	for rows.Next() {
		a := &Author{}
		check(rows.Scan(&a.ID, &a.Name))
		as = append(as, a)
	}
	check(rows.Err())
	return as
}

func (p *PQ) RecentArticles(n int) []*Article {
	var as []*Article
	rows, err := p.db.Query(`SELECT id, title, body, published_at FROM articles ORDER BY published_at DESC LIMIT $1`, n)
	check(err)
	defer rows.Close()
	for rows.Next() {
		a := &Article{}
		check(rows.Scan(&a.ID, &a.Title, &a.Body, &a.PublishedAt))
		as = append(as, a)
	}
	check(rows.Err())
	return as
}

var _ Implementation = &PQ{}
