package dbeval

import (
	"log"
	"time"

	"github.com/gocraft/dbr"
)

type DBR struct {
	db *dbr.Session
}

type dbrLogger struct{}

// Event receives a simple notification when various events occur
func (n *dbrLogger) Event(eventName string) {
	log.Println(eventName)
}

// EventKv receives a notification when various events occur along with
// optional key/value data
func (n *dbrLogger) EventKv(eventName string, kvs map[string]string) {
	log.Printf("%s %v", eventName, kvs)
}

// EventErr receives a notification of an error if one occurs
func (n *dbrLogger) EventErr(eventName string, err error) error {
	log.Printf("%s %v", eventName, err)
	return err
}

// EventErrKv receives a notification of an error if one occurs along with
// optional key/value data
func (n *dbrLogger) EventErrKv(eventName string, err error, kvs map[string]string) error {
	log.Printf("%s %v %v", eventName, err, err)
	return err
}

// Timing receives the time an event took to happen
func (n *dbrLogger) Timing(eventName string, nanoseconds int64) {
	log.Printf("%s %d", eventName, nanoseconds)
}

// TimingKv receives the time an event took to happen along with optional key/value data
func (n *dbrLogger) TimingKv(eventName string, nanoseconds int64, kvs map[string]string) {
	log.Printf("%s %v %v", eventName, nanoseconds, kvs)
}

func (p *DBR) Connect(ds string) {
	if p.db != nil {
		p.db.Close()
		p.db = nil
	}
	conn, err := dbr.Open("postgres", ds, nil)
	check(err)
	//p.db = conn.NewSession(&dbrLogger{})
	p.db = conn.NewSession(nil)
}

func (p *DBR) CreateDatabase() {
	_, err := p.db.Exec(createdb)
	check(err)
}

func (p *DBR) DropDatabase() {
	_, err := p.db.Exec(dropdb)
	check(err)
}

func (p *DBR) CreateSchema() {
	_, err := p.db.Exec(schema)
	check(err)
}

func (p *DBR) InsertAuthors(as []*Author) {
	tx, err := p.db.Begin()
	check(err)
	defer tx.RollbackUnlessCommitted()
	ib := tx.InsertInto("authors").Columns("id", "name")
	for _, a := range as {
		ib.Record(a)
	}
	_, err = ib.Exec()
	check(err)
	check(tx.Commit())
}

func (p *DBR) InsertArticles(as []*Article) {
	tx, err := p.db.Begin()
	check(err)
	defer tx.RollbackUnlessCommitted()
	ib := tx.InsertInto("articles").Columns("id", "title", "body", "published_at")
	for _, a := range as {
		// bug in time encoding, drops tz on insert, so convert to local and insert as UTC
		t := a.PublishedAt.Local()
		a.PublishedAt = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC)
		ib.Record(a)
	}
	_, err = ib.Exec()
	check(err)
	check(tx.Commit())
}

func (p *DBR) FindAuthorByID(id int64) *Author {
	a := &Author{}
	check(p.db.Select("*").From("authors").Where("id = ? ", id).LoadOne(a))
	return a
}

func (p *DBR) FindAuthorsByName(name string) []*Author {
	var as []*Author
	_, err := p.db.Select("*").From("authors").Where("name = ?", name).Load(&as)
	check(err)
	return as
}

func (p *DBR) RecentArticles(n int) []*Article {
	var as []*Article
	_, err := p.db.Select("*").From("articles").OrderDir("published_at", false).Limit(uint64(n)).Load(&as)
	check(err)
	return as
}

var _ Implementation = &DBR{}
