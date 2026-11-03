package dbeval

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"sync/atomic"
	"testing"
	"time"

	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"
)

func TestImplementations(t *testing.T) {
	for _, impl := range impls() {
		implTest(t, impl)
	}
}

// poolSize is how many distinct fake values the insert benchmarks cycle through. The values are generated before
// the timer starts: fake's text generation is slow enough (~470µs for 10 article bodies) that generating inside the
// loop dominated the measurement.
const poolSize = 1000

func BenchmarkInsertAuthors(b *testing.B) {
	var next int64
	names := make([]string, poolSize)
	for i := range names {
		names[i] = fake.FullName()
	}

	for _, impl := range impls() {
		b.Run(fmt.Sprintf("%T", impl), func(b *testing.B) {
			withTestDB(impl, func() {
				b.ReportAllocs()
				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						authors := make([]*Author, 10)
						for i := range authors {
							id := atomic.AddInt64(&next, 1)
							authors[i] = &Author{ID: id, Name: names[id%poolSize]}
						}
						impl.InsertAuthors(authors)
					}
				})
			})
		})
	}

	fmt.Println()
}

func BenchmarkInsertArticles(b *testing.B) {
	var next int64
	pool := make([]Article, poolSize)
	for i := range pool {
		pool[i] = Article{
			Title:       fake.Title(),
			Body:        fake.EmailBody(),
			PublishedAt: time.Unix(rand.Int64N(time.Now().Unix()), 0),
		}
	}

	for _, impl := range impls() {
		b.Run(fmt.Sprintf("%T", impl), func(b *testing.B) {
			withTestDB(impl, func() {
				b.ReportAllocs()
				b.ResetTimer()

				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						articles := make([]*Article, 10)
						for i := range articles {
							id := atomic.AddInt64(&next, 1)
							a := pool[id%poolSize]
							a.ID = id
							articles[i] = &a
						}
						impl.InsertArticles(articles)
					}
				})
			})
		})
	}

	fmt.Println()
}

func BenchmarkFindAuthorByID(b *testing.B) {
	n := int64(10000)
	authors := make([]*Author, n)
	for i := range authors {
		authors[i] = &Author{ID: int64(i + 1), Name: fake.FullName()}
	}

	for _, impl := range impls() {
		withTestDB(impl, func() {
			impl.InsertAuthors(authors)
			b.Run(fmt.Sprintf("%T", impl), func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()

				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						impl.FindAuthorByID(rand.Int64N(n) + 1)
					}
				})
			})
		})
	}

	fmt.Println()
}

func BenchmarkFindAuthorByName(b *testing.B) {
	authorNames := make([]string, 10)
	for i := range authorNames {
		authorNames[i] = fake.FullName()
	}
	n := int64(10000)
	authors := make([]*Author, n)
	for i := range authors {
		authors[i] = &Author{ID: int64(i + 1), Name: authorNames[rand.IntN(len(authorNames))]}
	}

	for _, impl := range impls() {
		withTestDB(impl, func() {
			impl.InsertAuthors(authors)
			b.Run(fmt.Sprintf("%T", impl), func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()

				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						impl.FindAuthorsByName(authorNames[rand.IntN(len(authorNames))])
					}
				})
			})
		})
	}

	fmt.Println()
}

func BenchmarkRecentArticles(b *testing.B) {
	articles := make([]*Article, 10000)
	for i := range articles {
		articles[i] = &Article{
			ID:          int64(i + 1),
			Title:       fake.Title(),
			Body:        fake.EmailBody(),
			PublishedAt: time.Unix(rand.Int64N(time.Now().Unix()), 0),
		}
	}
	for _, impl := range impls() {
		withTestDB(impl, func() {
			impl.InsertArticles(articles)
			b.Run(fmt.Sprintf("%T", impl), func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						impl.RecentArticles(2000)
					}
				})
			})
		})
	}

	fmt.Println()
}

func impls() []Implementation {
	impls := []Implementation{
		&Memory{},
		&PQ{},
		&PGX{},
		&PGXStdlib{},
		&UpperDB{},
		&SQLX{},
		&Gorm{},
		&DBR{},
		&GoPG{},
		&Xorm{},
		&Ozzo{},
		&Bun{},
		&BunPGX{},
		&Ent{},
		&SQLC{},
	}
	rand.Shuffle(len(impls), func(i, j int) {
		impls[i], impls[j] = impls[j], impls[i]
	})
	return impls
}

func withTestDB(impl Implementation, f func()) {
	impl.Connect("sslmode=disable dbname=postgres", 0, 0, 1)
	impl.DropDatabase()
	impl.CreateDatabase()
	impl.Connect("sslmode=disable dbname="+testDatabaseName, 0, 40, 40)
	impl.CreateSchema()
	defer func() {
		impl.Connect("sslmode=disable dbname=postgres", 0, 0, 1)
		impl.DropDatabase()
	}()
	f()
}

func implTest(t *testing.T, impl Implementation) {
	t.Run(fmt.Sprintf("%T", impl), func(t *testing.T) {
		withTestDB(impl, func() {

			impl.InsertAuthors([]*Author{
				{
					ID:   1234,
					Name: "Dave",
				},
				{
					ID:   2345,
					Name: "Not Dave",
				},
				{
					ID:   3456,
					Name: "Dave",
				},
			})
			assert.Equal(t, &Author{ID: 1234, Name: "Dave"}, impl.FindAuthorByID(1234))
			assert.Equal(t, &Author{ID: 2345, Name: "Not Dave"}, impl.FindAuthorByID(2345))
			assert.Equal(t, []*Author{{ID: 2345, Name: "Not Dave"}}, impl.FindAuthorsByName("Not Dave"))
			assert.Equal(t, []*Author{{ID: 1234, Name: "Dave"}, {ID: 3456, Name: "Dave"}}, impl.FindAuthorsByName("Dave"))

			impl.InsertArticles([]*Article{
				{
					ID:          123,
					PublishedAt: time.Date(2012, 1, 1, 0, 0, 0, 0, time.Local),
				},
				{
					ID:          234,
					PublishedAt: time.Date(2010, 1, 1, 0, 0, 0, 0, time.Local),
				},
				{
					ID:          345,
					PublishedAt: time.Date(2013, 1, 1, 0, 0, 0, 0, time.Local),
				},
			})

			expected, err := json.Marshal([]*Article{
				{
					ID:          345,
					PublishedAt: time.Date(2013, 1, 1, 0, 0, 0, 0, time.Local),
				},
				{
					ID:          123,
					PublishedAt: time.Date(2012, 1, 1, 0, 0, 0, 0, time.Local),
				},
			})
			assert.NoError(t, err)

			// Compare instants: drivers may return times in UTC rather than time.Local.
			recent := impl.RecentArticles(2)
			for _, a := range recent {
				a.PublishedAt = a.PublishedAt.In(time.Local)
			}
			actual, err := json.Marshal(recent)
			assert.NoError(t, err)

			assert.Equal(t, string(expected), string(actual))
		})
	})
}
