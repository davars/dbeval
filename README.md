# dbeval
Comparison of Go libraries for working with PostgreSQL

Runs benchmarks on the same "typical" database workflow implemented in various ways.

#### [Latest Results](https://github.com/davars/dbeval/blob/master/results.txt)

## Warning: calls `DROP DATABASE` on local postgres instance! 

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

### Implementated (suggestions welcome)
- [x] [github.com/lib/pq](https://godoc.org/github.com/lib/pq)
- [x] [github.com/jackc/pgx](https://godoc.org/github.com/jackc/pgx)
- [x] [github.com/jackc/pgx/stdlib](https://godoc.org/github.com/jackc/pgx/stdlib)
- [x] [upper.io/db.v3](https://godoc.org/upper.io/db.v3)
- [x] [github.com/jmoiron/sqlx](https://godoc.org/github.com/jmoiron/sqlx)
- [x] [github.com/jinzhu/gorm](https://godoc.org/github.com/jinzhu/gorm)
- [x] [github.com/gobuffalo/pop](https://godoc.org/github.com/gobuffalo/pop)
- [x] [github.com/gocraft/dbr](https://godoc.org/github.com/gocraft/dbr)
- [ ] ~~[github.com/astaxie/beego/orm](https://godoc.org/github.com/astaxie/beego/orm)~~ Can't figure out how to disconnect in order to drop test database
- [x] [github.com/go-xorm/xorm](https://godoc.org/github.com/go-xorm/xorm)
- [ ] [gopkg.in/gorp.v2](https://godoc.org/gopkg.in/gorp.v2)
- [ ] [github.com/jmoiron/modl](https://godoc.org/github.com/jmoiron/modl)
- [ ] [github.com/volatiletech/sqlboiler](https://godoc.org/github.com/volatiletech/sqlboiler)
- [ ] [gopkg.in/src-d/go-kallax.v1](https://godoc.org/gopkg.in/src-d/go-kallax.v1)
- [x] [github.com/go-pg/pg](https://godoc.org/github.com/go-pg/pg)
