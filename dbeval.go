package dbeval

import (
	_ "embed"
	"net"
	"os"
	"strings"
	"time"

	"github.com/uptrace/bun"
)

type Author struct {
	bun.BaseModel `bun:"table:authors,alias:au" xorm:"-"`

	ID   int64  `db:"id" gorm:"primaryKey" xorm:"'id'" bun:"id,pk,autoincrement"`
	Name string `db:"name"`
}

func (Author) TableName() string {
	return "authors"
}

type Article struct {
	bun.BaseModel `bun:"table:articles,alias:ar" xorm:"-"`

	ID          int64     `db:"id" gorm:"primaryKey" xorm:"'id'" bun:"id,pk,autoincrement"`
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

//go:embed schema.sql
var schema string

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// pgAddr returns the network and address of the postgres server named by the
// PGHOST and PGPORT environment variables (defaulting to localhost:5432), for
// drivers that don't read them on their own.  A PGHOST starting with "/" is a
// unix socket directory.
func pgAddr() (network, addr string) {
	host, port := os.Getenv("PGHOST"), os.Getenv("PGPORT")
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	if strings.HasPrefix(host, "/") {
		return "unix", host + "/.s.PGSQL." + port
	}
	return "tcp", net.JoinHostPort(host, port)
}
