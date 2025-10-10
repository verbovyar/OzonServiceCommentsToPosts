package pubsub

import (
	"sync"

	"ozonProject/internal/models"
)

type Bus struct {
	mu   sync.RWMutex
	subs map[int64]map[chan *models.Comment]struct{}
}

func New() *Bus {
	return &Bus{subs: make(map[int64]map[chan *models.Comment]struct{})}
}

func (b *Bus) Subscribe(postID int64) chan *models.Comment {
	ch := make(chan *models.Comment, 1)
	b.mu.Lock()

	if _, ok := b.subs[postID]; !ok {
		b.subs[postID] = make(map[chan *models.Comment]struct{})
	}

	b.subs[postID][ch] = struct{}{}
	b.mu.Unlock()

	return ch
}

func (b *Bus) Unsubscribe(postID int64, ch chan *models.Comment) {
	b.mu.Lock()
	if m, ok := b.subs[postID]; ok {
		if _, ok := m[ch]; ok {
			delete(m, ch)
			close(ch)
		}
		if len(m) == 0 {
			delete(b.subs, postID)
		}
	}
	b.mu.Unlock()
}

func (b *Bus) Publish(c *models.Comment) {
	b.mu.RLock()
	m := b.subs[c.PostID]
	var targets []chan *models.Comment
	for ch := range m {
		targets = append(targets, ch)
	}
	b.mu.RUnlock()

	for _, ch := range targets {
		select {
		case ch <- c:
		default:
		}
	}
}
