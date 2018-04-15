package dbeval

import (
	"database/sql"

	"github.com/davars/dbeval/models"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

type Boiler struct {
	db *sql.DB
}

func (b *Boiler) Connect(ds string) {
	if b.db != nil {
		check(b.db.Close())
		b.db = nil
	}
	var err error
	b.db, err = sql.Open("pgx", ds)
	check(err)
}

func (b *Boiler) CreateDatabase() {
	_, err := b.db.Exec(createdb)
	check(err)
}

func (b *Boiler) DropDatabase() {
	_, err := b.db.Exec(dropdb)
	check(err)
}

func (b *Boiler) CreateSchema() {
	_, err := b.db.Exec(schema)
	check(err)
}

func (b *Boiler) InsertAuthors(as []*Author) {
	for i := range as {
		ma := &models.Author{ID: as[i].ID, Name: as[i].Name}
		check(ma.Insert(b.db))
	}
}

func (b *Boiler) InsertArticles(as []*Article) {
	for i := range as {
		ma := &models.Article{ID: as[i].ID, Title: as[i].Title, Body: as[i].Body, PublishedAt: as[i].PublishedAt}
		check(ma.Insert(b.db))
	}
}

func (b *Boiler) FindAuthorByID(id int64) *Author {
	ma, err := models.Authors(b.db, qm.Where("id=?", id)).One()
	check(err)
	return &Author{ID: ma.ID, Name: ma.Name}
}

func (b *Boiler) FindAuthorsByName(name string) []*Author {
	mas, err := models.Authors(b.db, qm.Where("name=?", name)).All()
	check(err)
	as := make([]*Author, len(mas))
	for i := range mas {
		as[i] = &Author{ID: mas[i].ID, Name: mas[i].Name}
	}
	return as
}

func (b *Boiler) RecentArticles(n int) []*Article {
	mas, err := models.Articles(b.db, qm.OrderBy("published_at DESC"), qm.Limit(n)).All()
	check(err)
	as := make([]*Article, len(mas))
	for i := range mas {
		ma := mas[i]
		as[i] = &Article{ID: ma.ID, Title: ma.Title, Body: ma.Body, PublishedAt: ma.PublishedAt}
	}
	return as
}

var _ Implementation = &Boiler{}
