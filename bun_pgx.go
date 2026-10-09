package dbeval

import (
	"database/sql"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

// Notes:
// - same as Bun, but on top of the pgx stdlib driver instead of bun's pgdriver

type BunPGX struct {
	Bun
}

func (b *BunPGX) Connect(ds string, connLifetime time.Duration, idleConns, openConns int) {
	if b.db != nil {
		check(b.db.Close())
		b.db = nil
	}

	sqldb, err := sql.Open("pgx", ds)
	check(err)

	b.db = bun.NewDB(sqldb, pgdialect.New())
	b.db.SetConnMaxLifetime(connLifetime)
	b.db.SetMaxIdleConns(idleConns)
	b.db.SetMaxOpenConns(openConns)
}

var _ Implementation = &BunPGX{}
