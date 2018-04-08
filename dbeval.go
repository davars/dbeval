package dbeval

import (
	"time"
)

type Author struct {
	ID   int64  `db:"id" gorm:"PRIMARY_KEY"`
	Name string `db:"name"`
}

type Article struct {
	ID          int64     `db:"id" gorm:"PRIMARY_KEY"`
	Title       string    `db:"title"`
	Body        string    `db:"body"`
	PublishedAt time.Time `db:"published_at"`
}

type Implementation interface {
	Connect(string)
	CreateDatabase()
	DropDatabase()
	CreateSchema()
	InsertAuthors([]*Author)
	InsertArticles([]*Article)
	FindAuthorByID(int64) *Author
	FindAuthorsByName(string) []*Author
	RecentArticles(int) []*Article
}

const dbname = "dbeval_db"

const createdb = `
CREATE DATABASE dbeval_db;
`

const dropdb = `
DROP DATABASE dbeval_db;
`

const schema = `
CREATE UNLOGGED TABLE authors (
	id bigserial not null,
	name text not null,
    PRIMARY KEY(id)
);

CREATE INDEX authors_by_name ON authors (name);

CREATE UNLOGGED TABLE articles (
	id bigserial not null,
	title text not null,
	body text not null,
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
