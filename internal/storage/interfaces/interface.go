package interfaces

import (
	"context"
	"ozonProject/internal/models"
)

type storeIface interface {
	CreatePost(ctx context.Context, title, content, author string, commentsEnabled bool) (*models.Post, error)
	GetPosts(ctx context.Context, limit, offset int) ([]*models.Post, error)
	GetPostByID(ctx context.Context, id int64) (*models.Post, error)

	CreateComment(ctx context.Context, postID int64, parentID *int64, author, content string) (*models.Comment, error)
	GetComments(ctx context.Context, postID int64, parentID *int64, limit, offset int) ([]*models.Comment, error)
	EnsureCommentsEnabled(ctx context.Context, postID int64) error
}
