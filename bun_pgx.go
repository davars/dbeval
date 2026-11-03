package dbeval

import (
	"database/sql"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

// Notes:
// - same as Bun, but on top of the pgx stdlib driver instead of bun's pgdriver

// bun interpolates arguments into the SQL text itself, so pgx's default statement cache can never hit for
// parameter-varying queries: each distinct query would cost an extra round trip to prepare. The simple protocol sends
// each query in one round trip, at the cost of text rather than binary result decoding on large reads.
const bunPGXParams = " default_query_exec_mode=simple_protocol"

type BunPGX struct {
	Bun
}

func (b *BunPGX) Connect(ds string, connLifetime time.Duration, idleConns, openConns int) {
	if b.db != nil {
		check(b.db.Close())
		b.db = nil
	}

	sqldb, err := sql.Open("pgx", ds+bunPGXParams)
	check(err)

	b.db = bun.NewDB(sqldb, pgdialect.New())
	b.db.SetConnMaxLifetime(connLifetime)
	b.db.SetMaxIdleConns(idleConns)
	b.db.SetMaxOpenConns(openConns)
}

var _ Implementation = &BunPGX{}
