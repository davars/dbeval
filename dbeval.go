package dbeval

import (
	"time"
)

type Author struct {
	ID   int64  `db:"id" gorm:"PRIMARY_KEY" xorm:"'id'"`
	Name string `db:"name"`
}

func (Author) TableName() string {
	return "authors"
}

type Article struct {
	ID          int64     `db:"id" gorm:"PRIMARY_KEY" xorm:"'id'"`
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
