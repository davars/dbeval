package dbeval

import (
	"github.com/uptrace/bun"
	"time"
)

type Author struct {
	bun.BaseModel `bun:"table:authors,alias:au" xorm:"-"`

	ID   int64  `db:"id" gorm:"PRIMARY_KEY" xorm:"'id'" bun:"id,pk,autoincrement"`
	Name string `db:"name"`
}

func (Author) TableName() string {
	return "authors"
}

type Article struct {
	bun.BaseModel `bun:"table:articles,alias:ar" xorm:"-"`

	ID          int64     `db:"id" gorm:"PRIMARY_KEY" xorm:"'id'" bun:"id,pk,autoincrement"`
	Title       string    `db:"title"`
	Body        string    `db:"body"`
	PublishedAt time.Time `db:"published_at"`
}

func (Article) TableName() string {
	return "articles"
}

type Implementation interface {
	Connect(string, time.Duration, int, int)
	CreateDatabase()
	DropDatabase()
	CreateSchema()
	InsertAuthors([]*Author)
	InsertArticles([]*Article)
	FindAuthorByID(int64) *Author
	FindAuthorsByName(string) []*Author
	RecentArticles(int) []*Article
}

const testDatabaseName = "dbeval_db"

const createdb = `
CREATE DATABASE ` + testDatabaseName + `;`

const dropdb = `
DROP DATABASE IF EXISTS ` + testDatabaseName + `;`

const schema = `
CREATE UNLOGGED TABLE authors (
	id bigserial not null,
	name text not null default '',
    PRIMARY KEY(id)
);

CREATE INDEX authors_by_name ON authors (name);

CREATE UNLOGGED TABLE articles (
	id bigserial not null,
	title text not null default '',
	body text not null default '',
	published_at timestamp with time zone not null, 
    PRIMARY KEY(id)
);

CREATE INDEX articles_by_published_at ON articles (published_at);
`

func check(err error) {
	if err != nil {
		panic(err)
	}
}
