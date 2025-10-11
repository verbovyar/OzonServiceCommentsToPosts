package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"ozonProject/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PostgresStorage struct {
	pool PgxPoolIface
}

func NewPostgresStorage(pool PgxPoolIface) *PostgresStorage {
	return &PostgresStorage{pool: pool}
}

func (s *PostgresStorage) CreatePost(ctx context.Context, title, content, author string, commentsEnabled bool) (*models.Post, error) {
	id := uuid.New().String()

	const query = `
		INSERT INTO posts (id, title, content, author, comments_enabled)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, title, content, author, comments_enabled, created_at
	`

	log.Printf("Create post query") //

	var p models.Post
	err := s.pool.QueryRow(ctx, query, id, title, content, author, commentsEnabled).Scan(
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

	log.Printf("Get post query.")

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

func (s *PostgresStorage) GetPostByID(ctx context.Context, id string) (*models.Post, error) {
	const query = `
		SELECT id, title, content, author, comments_enabled, created_at
		FROM posts
		WHERE id = $1
	`
	log.Printf("Get post by id query.")

	var p models.Post
	err := s.pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.Title, &p.Content, &p.Author, &p.CommentsEnabled, &p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *PostgresStorage) CreateComment(ctx context.Context, postID string, parentID string, author, content string) (*models.Comment, error) {
	if err := s.EnsureCommentsEnabled(ctx, postID); err != nil {
		return nil, err
	}

	if parentID != "" {
		const queryPostId = `SELECT post_id FROM comments WHERE id = $1`

		log.Printf("Create comment query.")

		var parentPostID string
		if err := s.pool.QueryRow(ctx, queryPostId, parentID).Scan(&parentPostID); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, errors.New("parent comment not found")
			}
			return nil, err
		}
		if parentPostID != postID {
			return nil, errors.New("parent belongs to another post")
		}
	}

	id := uuid.New().String()
	const queryInsertComment = `
		INSERT INTO comments (id, post_id, parent_id, author, content)
		VALUES ($1, $2, $3 , $4, $5)
		RETURNING id, post_id, parent_id, author, content, created_at
	`

	log.Printf("Create comment (insert comment): %s", queryInsertComment) //

	var c models.Comment
	err := s.pool.QueryRow(ctx, queryInsertComment, id, postID, parentID, author, content).Scan(
		&c.ID, &c.PostID, &c.ParentID, &c.Author, &c.Content, &c.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (s *PostgresStorage) GetComments(ctx context.Context, postID string, parentID string, limit, offset int) ([]*models.Comment, error) {
	query := `
		SELECT id, post_id, parent_id, author, content, created_at
		FROM comments
		WHERE post_id = $1 %s
		ORDER BY id ASC
		LIMIT $2 OFFSET $3
	`

	log.Printf("Get comments query.") //

	var rows pgx.Rows
	var err error
	if parentID == "" {
		rows, err = s.pool.Query(ctx, fmt.Sprintf(query, "AND parent_id = '' "), postID, limit, offset)
	} else {
		rows, err = s.pool.Query(ctx, fmt.Sprintf(query, "AND parent_id = $4"), postID, limit, offset, parentID)
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

func (s *PostgresStorage) EnsureCommentsEnabled(ctx context.Context, postID string) error {
	const query = `SELECT comments_enabled FROM posts WHERE id = $1`

	log.Printf("Enable comments query.")

	var enabled bool
	if err := s.pool.QueryRow(ctx, query, postID).Scan(&enabled); err != nil {
		return err
	}

	if !enabled {
		return errors.New("comments disabled")
	}

	return nil
}
