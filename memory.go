package dbeval

import (
	"sort"
	"sync"

	"time"

	_ "github.com/lib/pq"
)

type Memory struct {
	sync.RWMutex
	authors        map[int64]*Author
	authorsByName  map[string][]*Author
	articles       map[int64]*Article
	articlesByDate []*Article
}

func (m *Memory) Connect(ds string, connLifetime time.Duration, idleConns, openConns int) {
	// Noop
}

func (m *Memory) CreateDatabase() {
	m.Lock()
	defer m.Unlock()
	m.authors = map[int64]*Author{}
	m.authorsByName = map[string][]*Author{}
	m.articles = map[int64]*Article{}
	m.articlesByDate = nil
}

func (m *Memory) DropDatabase() {
	m.CreateDatabase()
}

func (m *Memory) CreateSchema() {
	// Noop
}

func (m *Memory) InsertAuthors(as []*Author) {
	m.Lock()
	defer m.Unlock()
	for _, a := range as {
		m.authors[a.ID] = a
		m.authorsByName[a.Name] = append(m.authorsByName[a.Name], a)
	}
}

func (m *Memory) InsertArticles(as []*Article) {
	m.Lock()
	defer m.Unlock()
	for _, a := range as {
		m.articles[a.ID] = a
	}
	m.articlesByDate = append(m.articlesByDate, as...)
	sort.Slice(m.articlesByDate, func(i, j int) bool {
		return m.articlesByDate[i].PublishedAt.After(m.articlesByDate[j].PublishedAt)
	})
}

func (m *Memory) FindAuthorByID(id int64) *Author {
	m.RLock()
	defer m.RUnlock()
	return m.authors[id]
}

func (m *Memory) FindAuthorsByName(name string) []*Author {
	m.RLock()
	defer m.RUnlock()
	return m.authorsByName[name]
}

func (m *Memory) RecentArticles(n int) []*Article {
	m.RLock()
	defer m.RUnlock()
	return m.articlesByDate[0:n]
}

var _ Implementation = &Memory{}
