package dbeval

import (
	"context"
	"database/sql"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"

	"github.com/davars/dbeval/ent"
	"github.com/davars/dbeval/ent/article"
	"github.com/davars/dbeval/ent/author"
)

type Ent struct {
	db     *sql.DB
	client *ent.Client
}

func (e *Ent) Connect(ds string, connLifetime time.Duration, idleConns, openConns int) {
	if e.db != nil {
		check(e.db.Close())
		e.db = nil
	}
	var err error
	e.db, err = sql.Open("pgx", ds)
	check(err)
	e.db.SetConnMaxLifetime(connLifetime)
	e.db.SetMaxIdleConns(idleConns)
	e.db.SetMaxOpenConns(openConns)

	drv := entsql.OpenDB(dialect.Postgres, e.db)
	e.client = ent.NewClient(ent.Driver(drv))
}

func (e *Ent) CreateDatabase() {
	_, err := e.db.Exec(createdb)
	check(err)
}

func (e *Ent) DropDatabase() {
	_, err := e.db.Exec(dropdb)
	check(err)
}

// CreateSchema uses the shared schema.sql rather than e.client.Schema.Create so that ent runs against the same
// UNLOGGED tables as every other implementation; ent's migrations would create logged tables, which are slower to
// write to.
func (e *Ent) CreateSchema() {
	_, err := e.db.Exec(schema)
	check(err)
}

func (e *Ent) InsertAuthors(as []*Author) {
	_, err := e.client.Author.MapCreateBulk(as, func(a *ent.AuthorCreate, i int) {
		a.SetName(as[i].Name).SetID(as[i].ID)
	}).Save(context.Background())
	check(err)
}

func (e *Ent) InsertArticles(as []*Article) {
	_, err := e.client.Article.MapCreateBulk(as, func(a *ent.ArticleCreate, i int) {
		a.SetTitle(as[i].Title).SetBody(as[i].Body).SetID(as[i].ID).SetPublishedAt(as[i].PublishedAt)
	}).Save(context.Background())
	check(err)
}

func (e *Ent) FindAuthorByID(id int64) *Author {
	a, err := e.client.Author.Get(context.Background(), id)
	check(err)
	return &Author{
		ID:   a.ID,
		Name: a.Name,
	}
}

func (e *Ent) FindAuthorsByName(name string) []*Author {
	aes, err := e.client.Author.Query().Where(author.Name(name)).All(context.Background())
	check(err)
	as := make([]*Author, len(aes))
	for i := range aes {
		as[i] = &Author{
			ID:   aes[i].ID,
			Name: aes[i].Name,
		}
	}
	return as
}

func (e *Ent) RecentArticles(n int) []*Article {
	aes, err := e.client.Article.Query().Limit(n).Order(article.ByPublishedAt(entsql.OrderDesc())).All(context.Background())
	check(err)
	as := make([]*Article, len(aes))
	for i := range aes {
		as[i] = &Article{
			ID:          aes[i].ID,
			Title:       aes[i].Title,
			Body:        aes[i].Body,
			PublishedAt: aes[i].PublishedAt,
		}
	}
	return as
}

var _ Implementation = &Ent{}
