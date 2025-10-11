package storage

import (
	"context"
	"ozonProject/internal/models"

	"github.com/jackc/pgx/v5"
)

type PgxPoolIface interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type Storage interface {
	CreatePost(ctx context.Context, title, content, author string, commentsEnabled bool) (*models.Post, error)
	GetPosts(ctx context.Context, limit, offset int) ([]*models.Post, error)
	GetPostByID(ctx context.Context, id string) (*models.Post, error)

	CreateComment(ctx context.Context, postID string, parentID string, author, content string) (*models.Comment, error)
	GetComments(ctx context.Context, postID string, parentID string, limit, offset int) ([]*models.Comment, error)
	EnsureCommentsEnabled(ctx context.Context, postID string) error
}
