package dbeval

import (
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

// Notes:
// - NamedExec with a slice expands into a single multi-row INSERT, which is atomic on its own

type SQLX struct {
	db *sqlx.DB
}

func (p *SQLX) Connect(ds string, connLifetime time.Duration, idleConns, openConns int) {
	if p.db != nil {
		check(p.db.Close())
		p.db = nil
	}
	var err error
	p.db, err = sqlx.Open("pgx", ds)
	check(err)
	p.db.SetConnMaxLifetime(connLifetime)
	p.db.SetMaxIdleConns(idleConns)
	p.db.SetMaxOpenConns(openConns)
}

func (p *SQLX) CreateDatabase() {
	_, err := p.db.Exec(createdb)
	check(err)
}

func (p *SQLX) DropDatabase() {
	_, err := p.db.Exec(dropdb)
	check(err)
}

func (p *SQLX) CreateSchema() {
	_, err := p.db.Exec(schema)
	check(err)
}

func (p *SQLX) InsertAuthors(as []*Author) {
	_, err := p.db.NamedExec(`INSERT INTO authors (id, name) VALUES (:id, :name)`, as)
	check(err)
}

func (p *SQLX) InsertArticles(as []*Article) {
	_, err := p.db.NamedExec(`INSERT INTO articles (id, title, body, published_at) VALUES (:id, :title, :body, :published_at)`, as)
	check(err)
}

func (p *SQLX) FindAuthorByID(id int64) *Author {
	a := &Author{}
	check(p.db.Get(a, `SELECT id, name FROM authors WHERE id = $1`, id))
	return a
}

func (p *SQLX) FindAuthorsByName(name string) []*Author {
	var as []*Author
	check(p.db.Select(&as, `SELECT id, name FROM authors WHERE name = $1`, name))
	return as
}

func (p *SQLX) RecentArticles(n int) []*Article {
	var as []*Article
	check(p.db.Select(&as, `SELECT id, title, body, published_at FROM articles ORDER BY published_at DESC LIMIT $1`, n))
	return as
}

var _ Implementation = &SQLX{}
