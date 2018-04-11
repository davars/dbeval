package dbeval

import (
	"time"

	"github.com/gocraft/dbr"
)

// Notes:
// - trivial to insert multiple records with a single statement
// - inserts tz-ignorant times into tz-aware columns

type DBR struct {
	db *dbr.Session
}

func (p *DBR) Connect(ds string) {
	if p.db != nil {
		p.db.Close()
		p.db = nil
	}
	conn, err := dbr.Open("postgres", ds, nil)
	check(err)
	p.db = conn.NewSession(nil)
}

func (p *DBR) CreateDatabase() {
	_, err := p.db.Exec(createdb)
	check(err)
}

func (p *DBR) DropDatabase() {
	_, err := p.db.Exec(dropdb)
	check(err)
}

func (p *DBR) CreateSchema() {
	_, err := p.db.Exec(schema)
	check(err)
}

func (p *DBR) InsertAuthors(as []*Author) {
	tx, err := p.db.Begin()
	check(err)
	defer tx.RollbackUnlessCommitted()
	ib := tx.InsertInto("authors").Columns("id", "name")
	for _, a := range as {
		ib.Record(a)
	}
	_, err = ib.Exec()
	check(err)
	check(tx.Commit())
}

func (p *DBR) InsertArticles(as []*Article) {
	tx, err := p.db.Begin()
	check(err)
	defer tx.RollbackUnlessCommitted()
	ib := tx.InsertInto("articles").Columns("id", "title", "body", "published_at")
	for _, a := range as {
		// bug in time encoding, drops tz on insert, so convert to local and insert as UTC.  PG will add back local tz,
		// hopefully the same one as this program.
		t := a.PublishedAt.Local()
		a.PublishedAt = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC)
		ib.Record(a)
	}
	_, err = ib.Exec()
	check(err)
	check(tx.Commit())
}

func (p *DBR) FindAuthorByID(id int64) *Author {
	a := &Author{}
	check(p.db.Select("*").From("authors").Where("id = ? ", id).LoadOne(a))
	return a
}

func (p *DBR) FindAuthorsByName(name string) []*Author {
	var as []*Author
	_, err := p.db.Select("*").From("authors").Where("name = ?", name).Load(&as)
	check(err)
	return as
}

func (p *DBR) RecentArticles(n int) []*Article {
	var as []*Article
	_, err := p.db.Select("*").From("articles").OrderDir("published_at", false).Limit(uint64(n)).Load(&as)
	check(err)
	return as
}

var _ Implementation = &DBR{}
