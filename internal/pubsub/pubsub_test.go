package pubsub_test

import (
	"ozonProject/internal/models"
	"ozonProject/internal/pubsub"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBus_Publish_Subscribe(t *testing.T) {
	b := pubsub.New()
	postID := int64(42)

	ch := b.Subscribe(postID)
	defer b.Unsubscribe(postID, ch)

	msg := &models.Comment{ID: 1, PostID: postID, Content: "hello"}
	b.Publish(msg)

	select {
	case got := <-ch:
		require.Equal(t, msg, got)
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timeout waiting for message")
	}
}

func TestBus_Unsubscribe_ClosesChannel(t *testing.T) {
	b := pubsub.New()
	postID := int64(1)
	ch := b.Subscribe(postID)
	b.Unsubscribe(postID, ch)

	select {
	case _, ok := <-ch:
		require.False(t, ok, "channel should be closed")
	default:
		t.Fatal("channel not closed after unsubscribe")
	}
}
