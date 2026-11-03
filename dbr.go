package dbeval

import (
	"time"

	"github.com/gocraft/dbr/v2"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Notes:
// - trivial to insert multiple records with a single statement

// dbr interpolates arguments into the SQL text itself, so pgx's default statement cache can never hit for
// parameter-varying queries: each distinct query would cost an extra round trip to prepare. The simple protocol sends
// each query in one round trip, at the cost of text rather than binary result decoding on large reads.
const dbrPGXParams = " default_query_exec_mode=simple_protocol"

type DBR struct {
	db *dbr.Session
}

func (p *DBR) Connect(ds string, connLifetime time.Duration, idleConns, openConns int) {
	if p.db != nil {
		p.db.Close()
		p.db = nil
	}
	conn, err := dbr.Open("pgx", ds+dbrPGXParams, nil)
	check(err)
	conn.SetConnMaxLifetime(connLifetime)
	conn.SetMaxIdleConns(idleConns)
	conn.SetMaxOpenConns(openConns)
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
	ib := p.db.InsertInto("authors").Columns("id", "name")
	for _, a := range as {
		ib.Record(a)
	}
	_, err := ib.Exec()
	check(err)
}

func (p *DBR) InsertArticles(as []*Article) {
	ib := p.db.InsertInto("articles").Columns("id", "title", "body", "published_at")
	for _, a := range as {
		ib.Record(a)
	}
	_, err := ib.Exec()
	check(err)
}

func (p *DBR) FindAuthorByID(id int64) *Author {
	a := &Author{}
	check(p.db.Select("*").From("authors").Where("id = ?", id).LoadOne(a))
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
