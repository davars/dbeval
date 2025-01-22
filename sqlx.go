package dbeval

import (
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

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
	p.db.SetConnMaxLifetime(connLifetime)
	p.db.SetMaxIdleConns(idleConns)
	p.db.SetMaxOpenConns(openConns)
	check(err)
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
	tx, err := p.db.Beginx()
	defer func() {
		if err == nil {
			check(tx.Commit())
		} else {
			check(tx.Rollback())
		}
	}()
	check(err)
	for _, a := range as {
		_, err = tx.NamedExec(`INSERT INTO authors (id, name) VALUES (:id, :name)`, a)
		check(err)
	}
}

func (p *SQLX) InsertArticles(as []*Article) {
	tx, err := p.db.Beginx()
	defer func() {
		if err == nil {
			check(tx.Commit())
		} else {
			check(tx.Rollback())
		}
	}()
	check(err)
	for _, a := range as {
		_, err = tx.NamedExec(`INSERT INTO articles (id, title, body, published_at) VALUES (:id, :title, :body, :published_at)`, a)
		check(err)
	}
}

func (p *SQLX) FindAuthorByID(id int64) *Author {
	a := &Author{}
	check(p.db.QueryRowx(`SELECT id, name FROM authors WHERE id = $1`, id).StructScan(a))
	return a
}

func (p *SQLX) FindAuthorsByName(name string) []*Author {
	var as []*Author
	rows, err := p.db.Queryx(`SELECT id, name FROM authors WHERE name = $1`, name)
	check(err)
	defer rows.Close()
	for rows.Next() {
		a := &Author{}
		err = rows.StructScan(a)
		check(err)
		as = append(as, a)
	}
	check(rows.Err())
	return as
}

func (p *SQLX) RecentArticles(n int) []*Article {
	var as []*Article
	rows, err := p.db.Queryx(`SELECT id, title, body, published_at FROM articles ORDER BY published_at DESC LIMIT $1`, n)
	check(err)
	defer rows.Close()
	for rows.Next() {
		a := &Article{}
		err = rows.StructScan(a)
		check(err)
		as = append(as, a)
	}
	check(rows.Err())
	return as
}

var _ Implementation = &SQLX{}
