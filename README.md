# dbeval
Comparison of Go database libraries

## Warning: calls `DROP DATABASE` on local postgres instance! 

```
# Run benchmarks
go test -v -bench=. github.com/davars/dbeval

# Reset if cleanup was not performed by previous run
psql postgres -c 'DROP DATABASE dbeval_db'
```
