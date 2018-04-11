package dbeval

import (
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/postgresql"
)

type UpperDB struct {
	db       sqlbuilder.Database
	authors  db.Collection
	articles db.Collection
}

func (p *UpperDB) Connect(ds string) {
	if p.db != nil {
		check(p.db.Close())
		p.db = nil
	}
	cs, err := postgresql.ParseURL(ds)
	check(err)
	p.db, err = postgresql.Open(cs)
	check(err)
}

func (p *UpperDB) CreateDatabase() {
	_, err := p.db.Exec(createdb)
	check(err)
}

func (p *UpperDB) DropDatabase() {
	_, err := p.db.Exec(dropdb)
	check(err)
}

func (p *UpperDB) CreateSchema() {
	_, err := p.db.Exec(schema)
	check(err)
	p.authors = p.db.Collection("authors")
	p.articles = p.db.Collection("articles")
}

func (p *UpperDB) InsertAuthors(as []*Author) {
	batch := p.db.InsertInto("authors").Batch(len(as))
	for _, a := range as {
		batch = batch.Values(a)
	}
	batch.Done()
	check(batch.Wait())
}

func (p *UpperDB) InsertArticles(as []*Article) {
	batch := p.db.InsertInto("articles").Batch(len(as))
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
