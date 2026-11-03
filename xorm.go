package dbeval

import (
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"xorm.io/xorm"
)

// Notes:
// - Insert with a slice issues a single multi-row INSERT

type Xorm struct {
	db *xorm.Engine
}

func (x *Xorm) Connect(ds string, connLifetime time.Duration, idleConns, openConns int) {
	if x.db != nil {
		check(x.db.Close())
		x.db = nil
	}
	var err error
	x.db, err = xorm.NewEngine("pgx", ds)
	check(err)
	x.db.SetConnMaxLifetime(connLifetime)
	x.db.SetMaxIdleConns(idleConns)
	x.db.SetMaxOpenConns(openConns)
}

func (x *Xorm) CreateDatabase() {
	_, err := x.db.Exec(createdb)
	check(err)
}

func (x *Xorm) DropDatabase() {
	_, err := x.db.Exec(dropdb)
	check(err)
}

func (x *Xorm) CreateSchema() {
	_, err := x.db.Exec(schema)
	check(err)
}

func (x *Xorm) InsertAuthors(as []*Author) {
	_, err := x.db.Insert(&as)
	check(err)
}

func (x *Xorm) InsertArticles(as []*Article) {
	_, err := x.db.Insert(&as)
	check(err)
}

func (x *Xorm) FindAuthorByID(id int64) *Author {
	a := &Author{}
	has, err := x.db.ID(id).Get(a)
	check(err)
	if !has {
		panic("author not found")
	}
	return a
}

func (x *Xorm) FindAuthorsByName(name string) []*Author {
	var as []*Author
	check(x.db.Where("name = ?", name).Find(&as))
	return as
}

func (x *Xorm) RecentArticles(n int) []*Article {
	var as []*Article
	check(x.db.Desc("published_at").Limit(n).Find(&as))
	return as
}

var _ Implementation = &Xorm{}
