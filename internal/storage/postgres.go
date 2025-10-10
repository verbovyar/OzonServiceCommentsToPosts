package storage

import (
	"context"
	"errors"
	"fmt"
	"ozonProject/internal/models"

	"github.com/jackc/pgx/v4"
)

type PostgresStorage struct {
	pool PgxPoolIface
}

func NewPostgresStorage(pool PgxPoolIface) *PostgresStorage {
	return &PostgresStorage{pool: pool}
}

func (s *PostgresStorage) CreatePost(ctx context.Context, title, content, author string, commentsEnabled bool) (*models.Post, error) {
	const query = `
		INSERT INTO posts (title, content, author, comments_enabled)
		VALUES ($1,$2,$3,$4)
		RETURNING id, title, content, author, comments_enabled, created_at
	`
	var p models.Post
	err := s.pool.QueryRow(ctx, query, title, content, author, commentsEnabled).Scan(
		&p.ID, &p.Title, &p.Content, &p.Author, &p.CommentsEnabled, &p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *PostgresStorage) GetPosts(ctx context.Context, limit, offset int) ([]*models.Post, error) {
	const query = `
		SELECT id, title, content, author, comments_enabled, created_at
		FROM posts
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := s.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Post
	for rows.Next() {
		var p models.Post
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Author, &p.CommentsEnabled, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		out = append(out, &p)
	}

	return out, rows.Err()
}

func (s *PostgresStorage) GetPostByID(ctx context.Context, id int64) (*models.Post, error) {
	const query = `
		SELECT id, title, content, author, comments_enabled, created_at
		FROM posts
		WHERE id = $1
	`
	var p models.Post
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.Title, &p.Content, &p.Author, &p.CommentsEnabled, &p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *PostgresStorage) CreateComment(ctx context.Context, postID int64, parentID *int64, author, content string) (*models.Comment, error) {
	if err := s.EnsureCommentsEnabled(ctx, postID); err != nil {
		return nil, err
	}

	if parentID != nil {
		const queryParent = `SELECT post_id FROM comments WHERE id = $1`
		var parentPostID int64
		if err := s.pool.QueryRow(ctx, queryParent, *parentID).Scan(&parentPostID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, errors.New("parent comment not found")
			}
			return nil, err
		}
		if parentPostID != postID {
			return nil, errors.New("parent belongs to another post")
		}
	}

	const queryInsert = `
		INSERT INTO comments (post_id, parent_id, author, content)
		VALUES ($1,$2,$3,$4)
		RETURNING id, post_id, parent_id, author, content, created_at
	`
	var c models.Comment
	err := s.pool.QueryRow(ctx, queryInsert, postID, parentID, author, content).Scan(
		&c.ID, &c.PostID, &c.ParentID, &c.Author, &c.Content, &c.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (s *PostgresStorage) GetComments(ctx context.Context, postID int64, parentID *int64, limit, offset int) ([]*models.Comment, error) {
	base := `
		SELECT id, post_id, parent_id, author, content, created_at
		FROM comments
		WHERE post_id = $1 %s
		ORDER BY id ASC
		LIMIT $2 OFFSET $3
	`
	var rows pgx.Rows
	var err error
	if parentID == nil {
		rows, err = s.pool.Query(ctx, fmt.Sprintf(base, "AND parent_id IS NULL"), postID, limit, offset)
	} else {
		rows, err = s.pool.Query(ctx, fmt.Sprintf(base, "AND parent_id = $4"), postID, limit, offset, *parentID)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*models.Comment
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.ParentID, &c.Author, &c.Content, &c.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, &c)
	}

	return out, rows.Err()
}

func (s *PostgresStorage) EnsureCommentsEnabled(ctx context.Context, postID int64) error {
	const query = `SELECT comments_enabled FROM posts WHERE id = $1`
	var enabled bool
	if err := s.pool.QueryRow(ctx, query, postID).Scan(&enabled); err != nil {
		return err
	}

	if !enabled {
		return errors.New("comments disabled")
	}

	return nil
}
