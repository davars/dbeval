package dbeval

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type PQ struct {
	db *sql.DB
}

func (p *PQ) Connect(ds string) {
	if p.db != nil {
		check(p.db.Close())
		p.db = nil
	}
	var err error
	p.db, err = sql.Open("postgres", ds)
	check(err)
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

func (p *PQ) InsertAuthors(as []*Author) {
	tx, err := p.db.Begin()
	defer func() {
		if err == nil {
			check(tx.Commit())
		} else {
			check(tx.Rollback())
		}
	}()
	check(err)
	for _, a := range as {
		_, err = tx.Exec(`INSERT INTO authors (id, name) VALUES ($1, $2)`, a.ID, a.Name)
		check(err)
	}
}

func (p *PQ) InsertArticles(as []*Article) {
	tx, err := p.db.Begin()
	defer func() {
		if err == nil {
			check(tx.Commit())
		} else {
			check(tx.Rollback())
		}
	}()
	check(err)
	for _, a := range as {
		_, err = tx.Exec(`INSERT INTO articles (id, title, body, published_at) VALUES ($1, $2, $3, $4)`, a.ID, a.Title, a.Body, a.PublishedAt)
		check(err)
	}
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
		err = rows.Scan(&a.ID, &a.Name)
		check(err)
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
		err = rows.Scan(&a.ID, &a.Title, &a.Body, &a.PublishedAt)
		check(err)
		as = append(as, a)
	}
	check(rows.Err())
	return as
}

var _ Implementation = &PQ{}
