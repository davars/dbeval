package dbeval

import (
	"database/sql"

	"time"

	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/postgresql"
)

type UpperDB struct {
	db       *sql.DB
	sess     sqlbuilder.Database
	authors  db.Collection
	articles db.Collection
}

func (p *UpperDB) Connect(ds string, connLifetime time.Duration, idleConns, openConns int) {
	if p.sess != nil {
		check(p.sess.Close())
		p.sess = nil
		check(p.db.Close())
		p.db = nil
	}
	var err error
	p.db, err = sql.Open("pgx", ds)
	check(err)
	p.db.SetConnMaxLifetime(connLifetime)
	p.db.SetMaxIdleConns(idleConns)
	p.db.SetMaxOpenConns(openConns)
	p.sess, err = postgresql.New(p.db)
	check(err)
}

func (p *UpperDB) CreateDatabase() {
	_, err := p.sess.Exec(createdb)
	check(err)
}

func (p *UpperDB) DropDatabase() {
	_, err := p.sess.Exec(dropdb)
	check(err)
}

func (p *UpperDB) CreateSchema() {
	_, err := p.sess.Exec(schema)
	check(err)
	p.authors = p.sess.Collection("authors")
	p.articles = p.sess.Collection("articles")
}

func (p *UpperDB) InsertAuthors(as []*Author) {
	batch := p.sess.InsertInto("authors").Batch(len(as))
	for _, a := range as {
		batch = batch.Values(a)
	}
	batch.Done()
	check(batch.Wait())
}

func (p *UpperDB) InsertArticles(as []*Article) {
	batch := p.sess.InsertInto("articles").Batch(len(as))
	for _, a := range as {
		batch = batch.Values(a)
	}
	batch.Done()
	check(batch.Wait())
}

func (p *UpperDB) FindAuthorByID(id int64) *Author {
	a := &Author{}
	check(p.authors.Find("id", id).One(a))
	return a
}

func (p *UpperDB) FindAuthorsByName(name string) []*Author {
	var as []*Author
	check(p.authors.Find("name", name).All(&as))
	return as
}

func (p *UpperDB) RecentArticles(n int) []*Article {
	var as []*Article
	check(p.articles.Find().OrderBy("published_at DESC").Limit(n).All(&as))
	return as
}

var _ Implementation = &UpperDB{}
