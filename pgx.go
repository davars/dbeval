package dbeval

import (
	"github.com/jackc/pgx"
)

type PGX struct {
	db *pgx.ConnPool
}

func (p *PGX) Connect(ds string) {
	if p.db != nil {
		p.db.Close()
		p.db = nil
	}
	cc, err := pgx.ParseConnectionString(ds)
	p.db, err = pgx.NewConnPool(pgx.ConnPoolConfig{ConnConfig: cc})
	check(err)
}

func (p *PGX) CreateDatabase() {
	_, err := p.db.Exec(createdb)
	check(err)
}

func (p *PGX) DropDatabase() {
	_, err := p.db.Exec(dropdb)
	check(err)
}

func (p *PGX) CreateSchema() {
	_, err := p.db.Exec(schema)
	check(err)
}

func (p *PGX) InsertAuthors(as []*Author) {
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

func (p *PGX) InsertArticles(as []*Article) {
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
		_, err = tx.Exec("INSERT INTO articles (id, title, body, published_at) VALUES ($1, $2, $3, $4)", a.ID, a.Title, a.Body, a.PublishedAt)
		check(err)
	}
}

func (p *PGX) FindAuthorByID(id int64) *Author {
	a := &Author{}
	check(p.db.QueryRow(`SELECT id, name FROM authors WHERE id = $1`, id).Scan(&a.ID, &a.Name))
	return a
}

func (p *PGX) FindAuthorsByName(name string) []*Author {
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

func (p *PGX) RecentArticles(n int) []*Article {
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

var _ Implementation = &PGX{}