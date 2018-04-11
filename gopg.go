package dbeval

import (
	"os"

	"github.com/go-pg/pg"
	"upper.io/db.v3/postgresql"
)

// Notes:
// - trivial to insert multiple records with a single statement
// - treats zero value as DEFAULT rather than zero value 🤔
// - found a bug in my schema (not null without default)

type GoPG struct {
	db *pg.DB
}

func (g *GoPG) Connect(ds string) {
	if g.db != nil {
		check(g.db.Close())
		g.db = nil
	}
	url, err := postgresql.ParseURL(ds)
	check(err)
	g.db = pg.Connect(&pg.Options{
		User:     os.Getenv("USER"), // HACK
		Database: url.Database,
	})
}

func (g *GoPG) CreateDatabase() {
	_, err := g.db.Exec(createdb)
	check(err)
}

func (g *GoPG) DropDatabase() {
	_, err := g.db.Exec(dropdb)
	check(err)
}

func (g *GoPG) CreateSchema() {
	_, err := g.db.Exec(schema)
	check(err)
}

func (g *GoPG) InsertAuthors(as []*Author) {
	check(g.db.RunInTransaction(func(tx *pg.Tx) error {
		return g.db.Insert(&as)
	}))
}

func (g *GoPG) InsertArticles(as []*Article) {
	check(g.db.RunInTransaction(func(tx *pg.Tx) error {
		return g.db.Insert(&as)
	}))
}

func (g *GoPG) FindAuthorByID(id int64) *Author {
	a := &Author{ID: id}
	check(g.db.Select(a))
	return a
}

func (g *GoPG) FindAuthorsByName(name string) []*Author {
	var as []*Author
	check(g.db.Model(&as).Where("name = ?", name).Select())
	return as
}

func (g *GoPG) RecentArticles(n int) []*Article {
	var as []*Article
	check(g.db.Model(&as).Order("published_at DESC").Limit(n).Select())
	return as
}

var _ Implementation = &GoPG{}
