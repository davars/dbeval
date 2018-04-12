package dbeval

import (
	"github.com/go-xorm/xorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Xorm struct {
	db *xorm.Engine
}

func (x *Xorm) Connect(ds string) {
	if x.db != nil {
		check(x.db.Close())
		x.db = nil
	}
	var err error
	x.db, err = xorm.NewEngine("postgres", ds)
	check(err)
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
	_, err := x.db.Table("authors").Where("id = ?", id).Get(a)
	check(err)
	return a
}

func (x *Xorm) FindAuthorsByName(name string) []*Author {
	var as []*Author
	check(x.db.Table("authors").Where("name = ?", name).Find(&as))
	return as
}

func (x *Xorm) RecentArticles(n int) []*Article {
	var as []*Article
	check(x.db.Table("articles").Desc("published_at").Limit(n).Find(&as))
	return as
}

var _ Implementation = &Xorm{}
