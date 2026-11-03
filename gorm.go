package dbeval

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Notes:
// - can't use zero value as primary key
// - Create with a slice issues a single multi-row INSERT
// - SkipDefaultTransaction and PrepareStmt are the settings gorm's performance guide recommends

type Gorm struct {
	db *gorm.DB
}

func (g *Gorm) Connect(ds string, connLifetime time.Duration, idleConns, openConns int) {
	if g.db != nil {
		sqlDB, err := g.db.DB()
		check(err)
		check(sqlDB.Close())
		g.db = nil
	}
	var err error
	g.db, err = gorm.Open(postgres.Open(ds), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	check(err)
	sqlDB, err := g.db.DB()
	check(err)
	sqlDB.SetConnMaxLifetime(connLifetime)
	sqlDB.SetMaxIdleConns(idleConns)
	sqlDB.SetMaxOpenConns(openConns)
}

// execRaw runs harness DDL on the underlying *sql.DB: with PrepareStmt enabled, gorm would try to prepare the
// multi-statement schema, which postgres rejects.
func (g *Gorm) execRaw(query string) {
	sqlDB, err := g.db.DB()
	check(err)
	_, err = sqlDB.Exec(query)
	check(err)
}

func (g *Gorm) CreateDatabase() {
	g.execRaw(createdb)
}

func (g *Gorm) DropDatabase() {
	g.execRaw(dropdb)
}

func (g *Gorm) CreateSchema() {
	g.execRaw(schema)
}

func (g *Gorm) InsertAuthors(as []*Author) {
	check(g.db.Create(&as).Error)
}

func (g *Gorm) InsertArticles(as []*Article) {
	check(g.db.Create(&as).Error)
}

func (g *Gorm) FindAuthorByID(id int64) *Author {
	a := &Author{}
	check(g.db.First(a, id).Error)
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
