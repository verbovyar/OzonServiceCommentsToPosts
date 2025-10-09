package databases_test

import (
	"context"
	"errors"
	"ozonProject/internal/models"
	"ozonProject/internal/service"
	"testing"

	"github.com/stretchr/testify/require"
)

type fakeStore struct {
	commentsEnabled bool
}

func (f *fakeStore) CreatePost(ctx context.Context, title, content, author string, ce bool) (*models.Post, error) {
	return &models.Post{ID: 1, Title: title, Content: content, Author: author, CommentsEnabled: ce}, nil
}
func (f *fakeStore) GetPosts(ctx context.Context, limit, offset int) ([]*models.Post, error) {
	return []*models.Post{{ID: 1, Title: "t"}}, nil
}
func (f *fakeStore) GetPostByID(ctx context.Context, id int64) (*models.Post, error) {
	return &models.Post{ID: id, CommentsEnabled: f.commentsEnabled}, nil
}
func (f *fakeStore) CreateComment(ctx context.Context, postID int64, parentID *int64, author, content string) (*models.Comment, error) {
	return &models.Comment{ID: 10, PostID: postID, ParentID: parentID, Author: author, Content: content}, nil
}
func (f *fakeStore) GetComments(ctx context.Context, postID int64, parentID *int64, limit, offset int) ([]*models.Comment, error) {
	return []*models.Comment{{ID: 11, PostID: postID, ParentID: parentID, Author: "bob", Content: "hi"}}, nil
}
func (f *fakeStore) EnsureCommentsEnabled(ctx context.Context, postID int64) error {
	if !f.commentsEnabled {
		return errors.New("comments disabled")
	}
	return nil
}

func TestCreateComment_TooLong(t *testing.T) {
	t.Parallel()
	s := service.New(&fakeStore{commentsEnabled: true})
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
	s := service.New(&fakeStore{commentsEnabled: true})
	_, err := s.CreateComment(context.Background(), "1", nil, "me", "")
	require.Error(t, err)
	require.Contains(t, err.Error(), "content is empty")
}

func TestCreateComment_Disabled(t *testing.T) {
	t.Parallel()
	s := service.New(&fakeStore{commentsEnabled: false})
	_, err := s.CreateComment(context.Background(), "1", nil, "me", "ok")
	require.Error(t, err)
	require.Equal(t, "comments are disabled for this post", err.Error())
}

func TestCreateComment_Ok(t *testing.T) {
	t.Parallel()
	s := service.New(&fakeStore{commentsEnabled: true})
	c, err := s.CreateComment(context.Background(), "1", nil, "me", "ok")
	require.NoError(t, err)
	require.Equal(t, int64(10), c.ID)
	require.Equal(t, int64(1), c.PostID)
	require.Nil(t, c.ParentID)
	require.Equal(t, "ok", c.Content)
}

func TestListPosts_Defaults(t *testing.T) {
	t.Parallel()
	s := service.New(&fakeStore{commentsEnabled: true})
	posts, err := s.ListPosts(context.Background(), 0, -5)
	require.NoError(t, err)
	require.Len(t, posts, 1)
}
