package dbeval

import (
	"time"

	"github.com/go-ozzo/ozzo-dbx"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Notes:
// - not accepting slice-of-pointer-to-struct is annoying
// - no bulk insert :(

type Ozzo struct {
	db *dbx.DB
}

func (p *Ozzo) Connect(ds string, connLifetime time.Duration, idleConns, openConns int) {
	if p.db != nil {
		check(p.db.Close())
		p.db = nil
	}
	var err error
	p.db, err = dbx.Open("pgx", ds)
	check(err)
	p.db.DB().SetConnMaxLifetime(connLifetime)
	p.db.DB().SetMaxIdleConns(idleConns)
	p.db.DB().SetMaxOpenConns(openConns)
	//p.db.LogFunc = func(format string, a ...interface{}) {
	//	log.Printf(format, a...)
	//}
}

func (p *Ozzo) CreateDatabase() {
	_, err := p.db.DB().Exec(createdb)
	check(err)
}

func (p *Ozzo) DropDatabase() {
	_, err := p.db.DB().Exec(dropdb)
	check(err)
}

func (p *Ozzo) CreateSchema() {
	_, err := p.db.DB().Exec(schema)
	check(err)
}

func (p *Ozzo) InsertAuthors(as []*Author) {
	err := p.db.Transactional(func(tx *dbx.Tx) error {
		for _, a := range as {
			err := p.db.Model(a).Insert()
			if err != nil {
				return err
			}
		}
		return nil
	})
	check(err)
}

func (p *Ozzo) InsertArticles(as []*Article) {
	err := p.db.Transactional(func(tx *dbx.Tx) error {
		for _, a := range as {
			err := p.db.Model(a).Insert()
			if err != nil {
				return err
			}
		}
		return nil
	})
	check(err)
}

func (p *Ozzo) FindAuthorByID(id int64) *Author {
	a := &Author{}
	check(p.db.Select().Model(id, a))
	return a
}

func (p *Ozzo) FindAuthorsByName(name string) []*Author {
	var res []Author
	check(p.db.Select().Where(dbx.HashExp{"name": name}).All(&res))
	as := make([]*Author, len(res))
	for i := range res {
		as[i] = &res[i]
	}
	return as
}

func (p *Ozzo) RecentArticles(n int) []*Article {
	var res []Article
	check(p.db.Select().OrderBy("published_at DESC").Limit(int64(n)).All(&res))
	as := make([]*Article, len(res))
	for i := range res {
		as[i] = &res[i]
	}
	return as
}

var _ Implementation = &Ozzo{}
