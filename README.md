# dbeval
Comparison of Go libraries for working with PostgreSQL

Runs benchmarks on the same "typical" database workflow implemented in various ways.

#### [Latest Results](https://github.com/davars/dbeval/blob/master/results.txt)

## ⚠️ Warning: calls `DROP DATABASE` on local postgres instance! ⚠️

```
# Run benchmarks
go test -v -bench=. github.com/davars/dbeval

# Reset if cleanup was not performed by previous run
psql postgres -c 'DROP DATABASE dbeval_db'
```

### Style Note
Don't copy the `check` function to your own code!  More often than not you want to return errors to your caller, with
optional context from your error site.  When you get ready to show the error to your user, either print it out or
display it, optionally translating into an understandable error message.  I'm using panics here to simplify the test
harness.  The DB manipulation implementations then conform to that.  This is not a typical use-case so it's not likely
to be the best approach for your problem.

### Implementated
- [x] [github.com/lib/pq](https://pkg.go.dev/github.com/lib/pq)
- [x] [github.com/jackc/v5/pgx](https://pkg.go.dev/github.com/jackc/pgx/v5)
- [x] [github.com/jackc/pgx/v5/stdlib](https://pkg.go.dev/github.com/jackc/pgx/v5/stdlib)
- [x] [github.com/upper/db/v4](https://pkg.go.dev/github.com/upper/db/v4)
- [x] [github.com/jmoiron/sqlx](https://pkg.go.dev/github.com/jmoiron/sqlx)
- [x] [github.com/jinzhu/gorm](https://pkg.go.dev/github.com/jinzhu/gorm)
- [x] [github.com/gocraft/dbr](https://pkg.go.dev/github.com/gocraft/dbr)
- [x] [github.com/go-xorm/xorm](https://pkg.go.dev/github.com/go-xorm/xorm)
- [x] [github.com/go-pg/pg](https://pkg.go.dev/github.com/go-pg/pg)
- [x] [github.com/go-ozzo/ozzo-dbx](https://pkg.go.dev/github.com/go-ozzo/ozzo-dbx)
- [x] [https://github.com/uptrace/bun](https://pkg.go.dev/github.com/uptrace/bun)
- [x] [https://sqlc.dev/](https://sqlc.dev/)

### PRs Welcome
- [ ] [github.com/astaxie/beego/orm](https://pkg.go.dev/github.com/astaxie/beego/orm)
- [ ] [gopkg.in/gorp.v2](https://pkg.go.dev/gopkg.in/gorp.v2)
- [ ] [github.com/jmoiron/modl](https://pkg.go.dev/github.com/jmoiron/modl)
- [ ] [gopkg.in/src-d/go-kallax.v1](https://pkg.go.dev/gopkg.in/src-d/go-kallax.v1)
- [ ] [github.com/gobuffalo/pop](https://pkg.go.dev/github.com/gobuffalo/pop)
- [ ] [github.com/volatiletech/sqlboiler](https://pkg.go.dev/github.com/volatiletech/sqlboiler)
