package dbeval

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestMemoryRecentArticlesConcurrent checks that RecentArticles always returns articles newest first, and that a
// returned slice isn't reordered by a later sort, while inserts run concurrently.
func TestMemoryRecentArticlesConcurrent(t *testing.T) {
	m := &Memory{}
	m.CreateDatabase()

	var next atomic.Int64
	var wg sync.WaitGroup
	stop := make(chan struct{})
	for range 2 {
		wg.Go(func() {
			for {
				select {
				case <-stop:
					return
				default:
				}
				id := next.Add(1)
				m.InsertArticles([]*Article{{ID: id, PublishedAt: time.Unix(id, 0)}})
			}
		})
	}

	deadline := time.Now().Add(500 * time.Millisecond)
	for time.Now().Before(deadline) {
		as := m.RecentArticles(1 << 30)
		snapshot := append([]*Article(nil), as...)
		for i := 1; i < len(as); i++ {
			if as[i].PublishedAt.After(as[i-1].PublishedAt) {
				close(stop)
				wg.Wait()
				t.Fatalf("articles out of order at %d: %v after %v", i, as[i].PublishedAt, as[i-1].PublishedAt)
			}
		}
		m.RecentArticles(1 << 30) // may sort again; must not affect as
		for i := range as {
			if as[i] != snapshot[i] {
				close(stop)
				wg.Wait()
				t.Fatalf("returned slice changed at %d after a later call", i)
			}
		}
	}
	close(stop)
	wg.Wait()
}
