package dbeval

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Notes:
// - can't use zero value as primary key

type Gorm struct {
	db *gorm.DB
}

func (g *Gorm) Connect(ds string) {
	if g.db != nil {
		check(g.db.Close())
		g.db = nil
	}
	var err error
	g.db, err = gorm.Open("postgres", ds)
	check(err)
}

func (g *Gorm) CreateDatabase() {
	check(g.db.Exec(createdb).Error)
}

func (g *Gorm) DropDatabase() {
	check(g.db.Exec(dropdb).Error)
}

func (g *Gorm) CreateSchema() {
	check(g.db.Exec(schema).Error)
}

func (g *Gorm) InsertAuthors(as []*Author) {
	db := g.db.Begin()
	var err error
	defer func() {
		if err == nil {
			check(db.Commit().Error)
		} else {
			check(db.Rollback().Error)
		}
	}()
	for _, a := range as {
		err = db.Create(a).Error
		check(err)
	}
}

func (g *Gorm) InsertArticles(as []*Article) {
	db := g.db.Begin()
	var err error
	defer func() {
		if err == nil {
			check(db.Commit().Error)
		} else {
			check(db.Rollback().Error)
		}
	}()
	for _, a := range as {
		err = db.Create(a).Error
		check(err)
	}
}

func (g *Gorm) FindAuthorByID(id int64) *Author {
	a := &Author{}
	check(g.db.Where("id = ?", id).First(a).Error)
	return a
}

func (g *Gorm) FindAuthorsByName(name string) []*Author {
	var as []*Author
	check(g.db.Where("name = ?", name).Find(&as).Error)
	return as
}

func (g *Gorm) RecentArticles(n int) []*Article {
	var as []*Article
	check(g.db.Order("published_at DESC").Limit(n).Find(&as).Error)
	return as
}

var _ Implementation = &Gorm{}
