package storage_test

import (
	"context"
	"errors"
	"ozonProject/internal/models"
	"ozonProject/internal/service"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockStore struct {
	commentsEnabled bool
}

func (f *mockStore) CreatePost(ctx context.Context, title, content, author string, ce bool) (*models.Post, error) {
	return &models.Post{ID: "1", Title: title, Content: content, Author: author, CommentsEnabled: ce}, nil
}
func (f *mockStore) GetPosts(ctx context.Context, limit, offset int) ([]*models.Post, error) {
	return []*models.Post{{ID: "1", Title: "t"}}, nil
}
func (f *mockStore) GetPostByID(ctx context.Context, id string) (*models.Post, error) {
	return &models.Post{ID: id, CommentsEnabled: f.commentsEnabled}, nil
}
func (f *mockStore) CreateComment(ctx context.Context, postID string, parentID string, author, content string) (*models.Comment, error) {
	return &models.Comment{ID: "10", PostID: postID, ParentID: &parentID, Author: author, Content: content}, nil
}
func (f *mockStore) GetComments(ctx context.Context, postID string, parentID string, limit, offset int) ([]*models.Comment, error) {
	return []*models.Comment{{ID: "11", PostID: postID, ParentID: &parentID, Author: "bob", Content: "hi"}}, nil
}
func (f *mockStore) EnsureCommentsEnabled(ctx context.Context, postID string) error {
	if !f.commentsEnabled {
		return errors.New("comments disabled")
	}
	return nil
}

func TestCreateComment_TooLong(t *testing.T) {
	t.Parallel()
	s := service.New(&mockStore{commentsEnabled: true})
	body := make([]byte, 2001)
	for i := range body {
		body[i] = 'a'
	}
	_, err := s.CreateComment(context.Background(), "1", nil, "me", string(body))
	require.Error(t, err)
	require.Contains(t, err.Error(), "content too long")
}

func TestCreateComment_Empty(t *testing.T) {
	t.Parallel()
	s := service.New(&mockStore{commentsEnabled: true})
	_, err := s.CreateComment(context.Background(), "1", nil, "me", "")
	require.Error(t, err)
	require.Contains(t, err.Error(), "content is empty")
}

func TestCreateComment_Disabled(t *testing.T) {
	t.Parallel()
	s := service.New(&mockStore{commentsEnabled: false})
	_, err := s.CreateComment(context.Background(), "1", nil, "me", "ok")
	require.Error(t, err)
	require.Equal(t, "comments are disabled for this post", err.Error())
}

func TestCreateComment_Ok(t *testing.T) {
	t.Parallel()
	s := service.New(&mockStore{commentsEnabled: true})
	c, err := s.CreateComment(context.Background(), "1", nil, "me", "ok")
	require.NoError(t, err)
	require.Equal(t, "10", c.ID)
	require.Equal(t, "1", c.PostID)
	require.Equal(t, "", *c.ParentID)
	require.Equal(t, "ok", c.Content)
}

func TestListPosts_Defaults(t *testing.T) {
	t.Parallel()
	s := service.New(&mockStore{commentsEnabled: true})
	limit := 0
	offset := -5
	posts, err := s.ListPosts(context.Background(), &limit, &offset)
	require.NoError(t, err)
	require.Len(t, posts, 1)
}
