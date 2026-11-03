package dbeval

import (
	"slices"
	"sync"
	"time"
)

type Memory struct {
	sync.RWMutex
	authors        map[int64]*Author
	authorsByName  map[string][]*Author
	articles       map[int64]*Article
	articlesByDate []*Article // newest first once sorted
	unsorted       bool       // articlesByDate has appends since it was last sorted
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
	m.unsorted = false
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
	// Sort lazily on the next read rather than on every insert.
	m.articlesByDate = append(m.articlesByDate, as...)
	m.unsorted = true
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
	if m.unsorted {
		m.RUnlock()
		m.Lock()
		if m.unsorted {
			slices.SortStableFunc(m.articlesByDate, func(a, b *Article) int {
				return b.PublishedAt.Compare(a.PublishedAt)
			})
			m.unsorted = false
		}
		m.Unlock()
		m.RLock()
	}
	defer m.RUnlock()
	return m.articlesByDate[:min(n, len(m.articlesByDate))]
}

var _ Implementation = &Memory{}
