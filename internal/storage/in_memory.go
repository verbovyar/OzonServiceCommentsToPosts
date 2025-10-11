package storage

import (
	"context"
	"errors"
	"ozonProject/internal/models"
	"sync"
	"time"

	"github.com/google/uuid"
)

type postsStore struct {
	mu    sync.RWMutex
	byID  map[string]*models.Post
	order []string
}

func newPostsStore() *postsStore {
	return &postsStore{
		byID:  make(map[string]*models.Post),
		order: make([]string, 0, 200),
	}
}

func (s *postsStore) create(title, content, author string, commentsEnabled bool) *models.Post {
	s.mu.Lock()
	defer s.mu.Unlock()

	p := &models.Post{
		ID:              uuid.New().String(),
		Title:           title,
		Content:         content,
		Author:          author,
		CommentsEnabled: commentsEnabled,
		CreatedAt:       time.Now().UTC(),
	}

	s.byID[p.ID] = p
	s.order = append(s.order, p.ID)

	return p
}

func (s *postsStore) getByID(id string) (*models.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	p, ok := s.byID[id]
	if !ok {
		return nil, errors.New("post not found")
	}
	cp := *p

	return &cp, nil
}

func (s *postsStore) list(limit, offset int) ([]*models.Post, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	n := len(s.order)
	if offset >= n {
		return []*models.Post{}, nil
	}

	start := n - 1 - offset
	res := make([]*models.Post, 0, limit)
	for i := start; i >= 0 && len(res) < limit; i-- {
		id := s.order[i]
		if p, ok := s.byID[id]; ok {
			cp := *p
			res = append(res, &cp)
		}
	}

	return res, nil
}

type commentsStore struct {
	mu         sync.RWMutex
	byID       map[string]*models.Comment
	byPostRoot map[string][]string
	byParent   map[parentKey][]string
}

type parentKey struct {
	postID string
	parent string
}

func newCommentsStore() *commentsStore {
	return &commentsStore{
		byID:       make(map[string]*models.Comment),
		byPostRoot: make(map[string][]string),
		byParent:   make(map[parentKey][]string),
	}
}

func (s *commentsStore) create(postID string, parentID string, author, content string) *models.Comment {
	s.mu.Lock()
	defer s.mu.Unlock()

	c := &models.Comment{
		ID:        uuid.New().String(),
		PostID:    postID,
		ParentID:  &parentID,
		Author:    author,
		Content:   content,
		CreatedAt: time.Now().UTC(),
	}
	s.byID[c.ID] = c
	s.byPostRoot[postID] = append(s.byPostRoot[postID], c.ID)

	pk := parentKey{
		postID: postID,
	}

	if parentID != "" {
		pk.parent = parentID
	}
	s.byParent[pk] = append(s.byParent[pk], c.ID)

	return c
}

func (s *commentsStore) list(postID string, parentID string, limit, offset int) ([]*models.Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	pk := parentKey{
		postID: postID,
	}

	if parentID != "" {
		pk.parent = parentID
	}
	ids := s.byParent[pk]

	if offset >= len(ids) {
		return []*models.Comment{}, nil
	}

	end := offset + limit
	if end > len(ids) {
		end = len(ids)
	}
	slice := ids[offset:end]
	out := make([]*models.Comment, 0, len(slice))

	for _, id := range slice {
		if c, ok := s.byID[id]; ok {
			cp := *c
			out = append(out, &cp)
		}
	}

	return out, nil
}

type InMemoryStorage struct {
	posts    *postsStore
	comments *commentsStore
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		posts:    newPostsStore(),
		comments: newCommentsStore(),
	}
}

func (r *InMemoryStorage) CreatePost(ctx context.Context, title, content, author string, commentsEnabled bool) (*models.Post, error) {
	return r.posts.create(title, content, author, commentsEnabled), nil
}

func (r *InMemoryStorage) GetPosts(ctx context.Context, limit, offset int) ([]*models.Post, error) {
	return r.posts.list(limit, offset)
}

func (r *InMemoryStorage) GetPostByID(ctx context.Context, id string) (*models.Post, error) {
	return r.posts.getByID(id)
}

func (r *InMemoryStorage) CreateComment(ctx context.Context, postID string, parentID string, author, content string) (*models.Comment, error) {
	if _, err := r.posts.getByID(postID); err != nil {
		return nil, err
	}

	if parentID != "" {
		r.comments.mu.RLock()
		_, ok := r.comments.byID[parentID]
		r.comments.mu.RUnlock()
		if !ok {
			return nil, errors.New("parent comment not found")
		}
	}

	return r.comments.create(postID, parentID, author, content), nil
}

func (r *InMemoryStorage) GetComments(ctx context.Context, postID string, parentID string, limit, offset int) ([]*models.Comment, error) {
	return r.comments.list(postID, parentID, limit, offset)
}

func (r *InMemoryStorage) EnsureCommentsEnabled(ctx context.Context, postID string) error {
	p, err := r.posts.getByID(postID)
	if err != nil {
		return err
	}

	if !p.CommentsEnabled {
		return errors.New("comments disabled")
	}

	return nil
}
