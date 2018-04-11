package dbeval

import (
	"github.com/gobuffalo/pop"
	"upper.io/db.v3/postgresql"
)

// Notes:
// - can't use zero value as primary key
// - using raw SQL for insert because inserting with known IDs is not supported

type Pop struct {
	db *pop.Connection
}

func (g *Pop) Connect(ds string) {
	if g.db != nil {
		check(g.db.Close())
		g.db = nil
	}
	url, err := postgresql.ParseURL(ds)
	check(err)
	g.db, err = pop.NewConnection(&pop.ConnectionDetails{
		Dialect:  "postgres",
		Database: url.Database,
	})
	check(err)
	check(g.db.Open())
}

func (g *Pop) CreateDatabase() {
	check(g.db.RawQuery(createdb).Exec())
}

func (g *Pop) DropDatabase() {
	check(g.db.RawQuery(dropdb).Exec())
}

func (g *Pop) CreateSchema() {
	check(g.db.RawQuery(schema).Exec())
}

func (g *Pop) InsertAuthors(as []*Author) {
	err := g.db.Transaction(func(tx *pop.Connection) error {
		for _, a := range as {
			check(g.db.RawQuery(`INSERT INTO authors (id, name) VALUES ($1, $2)`, a.ID, a.Name).Exec())
		}
		return nil
	})
	check(err)
}

func (g *Pop) InsertArticles(as []*Article) {
	err := g.db.Transaction(func(tx *pop.Connection) error {
		for _, a := range as {
			check(g.db.RawQuery(`INSERT INTO articles (id, title, body, published_at) VALUES ($1, $2, $3, $4)`, a.ID, a.Title, a.Body, a.PublishedAt).Exec())
		}
		return nil
	})
	check(err)
}

func (g *Pop) FindAuthorByID(id int64) *Author {
	a := &Author{}
	check(g.db.Find(a, id))
	return a
}

func (g *Pop) FindAuthorsByName(name string) []*Author {
	var as []*Author
	check(g.db.Where("name = ?", name).All(&as))
	return as
}

func (g *Pop) RecentArticles(n int) []*Article {
	var as []*Article
	check(g.db.Order("published_at DESC").Limit(n).All(&as))
	return as
}

var _ Implementation = &Pop{}
