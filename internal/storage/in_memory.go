package storage

import (
	"context"
	"errors"
	"ozonProject/internal/models"
	"sync"
	"time"
)

type postsStore struct {
	mu    sync.RWMutex
	seq   int64
	byID  map[int64]*models.Post
	order []int64
}

func newPostsStore() *postsStore {
	return &postsStore{
		byID:  make(map[int64]*models.Post),
		order: make([]int64, 0, 200),
	}
}

func (s *postsStore) create(title, content, author string, commentsEnabled bool) *models.Post {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.seq++
	p := &models.Post{
		ID:              s.seq,
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

func (s *postsStore) getByID(id int64) (*models.Post, error) {
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
	seq        int64
	byID       map[int64]*models.Comment
	byPostRoot map[int64][]int64
	byParent   map[parentKey][]int64
}

type parentKey struct {
	postID int64
	parent int64
}

func newCommentsStore() *commentsStore {
	return &commentsStore{
		byID:       make(map[int64]*models.Comment),
		byPostRoot: make(map[int64][]int64),
		byParent:   make(map[parentKey][]int64),
	}
}

func (s *commentsStore) create(postID int64, parentID *int64, author, content string) *models.Comment {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.seq++
	c := &models.Comment{
		ID:        s.seq,
		PostID:    postID,
		ParentID:  parentID,
		Author:    author,
		Content:   content,
		CreatedAt: time.Now().UTC(),
	}
	s.byID[c.ID] = c
	s.byPostRoot[postID] = append(s.byPostRoot[postID], c.ID)

	pk := parentKey{
		postID: postID,
	}

	if parentID != nil {
		pk.parent = *parentID
	}
	s.byParent[pk] = append(s.byParent[pk], c.ID)

	return c
}

func (s *commentsStore) list(postID int64, parentID *int64, limit, offset int) ([]*models.Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	pk := parentKey{
		postID: postID,
	}

	if parentID != nil {
		pk.parent = *parentID
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

func (r *InMemoryStorage) GetPostByID(ctx context.Context, id int64) (*models.Post, error) {
	return r.posts.getByID(id)
}

func (r *InMemoryStorage) CreateComment(ctx context.Context, postID int64, parentID *int64, author, content string) (*models.Comment, error) {
	if _, err := r.posts.getByID(postID); err != nil {
		return nil, err
	}

	if parentID != nil {
		r.comments.mu.RLock()
		_, ok := r.comments.byID[*parentID]
		r.comments.mu.RUnlock()
		if !ok {
			return nil, errors.New("parent comment not found")
		}
	}

	return r.comments.create(postID, parentID, author, content), nil
}

func (r *InMemoryStorage) GetComments(ctx context.Context, postID int64, parentID *int64, limit, offset int) ([]*models.Comment, error) {
	return r.comments.list(postID, parentID, limit, offset)
}

func (r *InMemoryStorage) EnsureCommentsEnabled(ctx context.Context, postID int64) error {
	p, err := r.posts.getByID(postID)
	if err != nil {
		return err
	}

	if !p.CommentsEnabled {
		return errors.New("comments disabled")
	}

	return nil
}
