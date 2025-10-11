package storage_test

import (
	"context"
	"ozonProject/internal/storage"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/require"
)

func TestGetPosts_Ok(t *testing.T) {
	mockPool, err := pgxmock.NewConn()
	if err != nil {
		t.Fatal(err)
	}
	require.NoError(t, err)

	storage := storage.NewPostgresStorage(mockPool)

	firstTime := time.Now().UTC()
	secondTime := time.Now().UTC()

	rows := pgxmock.NewRows([]string{"id", "title", "content", "author", "comments_enabled", "created_at"}).
		AddRow("1", "first post", "Hello", "Yaroslav", true, firstTime).
		AddRow("2", "second post", "Hi", "Sergey", false, secondTime)

	const query = `
		SELECT id, title, content, author, comments_enabled, created_at
		FROM posts
		ORDER BY id DESC
		LIMIT \$1 OFFSET \$2
	`

	mockPool.ExpectQuery(query).WithArgs(10, 0).WillReturnRows(rows)
	posts, err := storage.GetPosts(context.Background(), 10, 0)

	require.NoError(t, err)
	require.Equal(t, "1", posts[0].ID)
	require.Equal(t, "second post", posts[1].Title)
	require.Equal(t, "Yaroslav", posts[0].Author)
}
