package dbeval

import (
	"context"
	"os"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/upper/db/v4/adapter/postgresql"
)

// Notes:
// - trivial to insert multiple records with a single statement
// - treats zero value as DEFAULT rather than zero value 🤔
// - found a bug in my schema (not null without default)

type GoPG struct {
	db *pg.DB
}

func (g *GoPG) Connect(ds string, connLifetime time.Duration, idleConns, openConns int) {
	if g.db != nil {
		check(g.db.Close())
		g.db = nil
	}
	url, err := postgresql.ParseURL(ds)
	check(err)
	network, addr := pgAddr()
	g.db = pg.Connect(&pg.Options{
		Network:  network,
		Addr:     addr,
		User:     os.Getenv("USER"), // HACK
		Database: url.Database,
		PoolSize: openConns,
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
	check(g.db.RunInTransaction(context.Background(), func(tx *pg.Tx) error {
		_, err := tx.Model(&as).Insert()
		return err
	}))
}

func (g *GoPG) InsertArticles(as []*Article) {
	check(g.db.RunInTransaction(context.Background(), func(tx *pg.Tx) error {
		_, err := tx.Model(&as).Insert()
		return err
	}))
}

func (g *GoPG) FindAuthorByID(id int64) *Author {
	a := &Author{ID: id}
	check(g.db.Model(a).WherePK().Select())
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
