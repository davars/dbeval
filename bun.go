package dbeval

import (
	"context"
	"database/sql"
	"os"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"upper.io/db.v3/postgresql"
)

// Notes:
// - trivial to insert multiple records with a single statement

type Bun struct {
	db *bun.DB
}

func (b *Bun) Connect(ds string, connLifetime time.Duration, idleConns, openConns int) {
	if b.db != nil {
		check(b.db.Close())
		b.db = nil
	}

	url, err := postgresql.ParseURL(ds)
	check(err)
	sqldb := sql.OpenDB(
		pgdriver.NewConnector(pgdriver.WithDatabase(url.Database), pgdriver.WithUser(os.Getenv("USER"))))

	b.db = bun.NewDB(sqldb, pgdialect.New())
	b.db.SetConnMaxLifetime(connLifetime)
	b.db.SetMaxIdleConns(idleConns)
	b.db.SetMaxOpenConns(openConns)
}

func (b *Bun) CreateDatabase() {
	_, err := b.db.Exec(createdb)
	check(err)
}

func (b *Bun) DropDatabase() {
	_, err := b.db.Exec(dropdb)
	check(err)
}

func (b *Bun) CreateSchema() {
	_, err := b.db.Exec(schema)
	check(err)
}

func (b *Bun) InsertAuthors(as []*Author) {
	check(b.db.RunInTx(context.Background(), nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().Model(&as).Exec(ctx)
		return err
	}))
}

func (b *Bun) InsertArticles(as []*Article) {
	check(b.db.RunInTx(context.Background(), nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().Model(&as).Exec(ctx)
		return err
	}))
}

func (b *Bun) FindAuthorByID(id int64) *Author {
	a := &Author{ID: id}
	check(b.db.NewSelect().Model(a).WherePK().Scan(context.Background()))
	return a
}

func (b *Bun) FindAuthorsByName(name string) []*Author {
	var as []*Author
	check(b.db.NewSelect().Model(&as).Where("name = ?", name).Scan(context.Background()))
	return as
}

func (b *Bun) RecentArticles(n int) []*Article {
	var as []*Article
	check(b.db.NewSelect().Model(&as).Order("published_at DESC").Limit(n).Scan(context.Background()))
	return as
}

var _ Implementation = &Bun{}
