# dbeval
Comparison of Go database libraries

Runs benchmarks on the same "typical" database workflow implemented in various ways.

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
