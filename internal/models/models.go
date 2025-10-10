package models

import "time"

type Post struct {
	ID              int64     `json:"id"`
	Title           string    `json:"title"`
	Content         string    `json:"content"`
	Author          string    `json:"author"`
	CommentsEnabled bool      `json:"commentsEnabled"`
	CreatedAt       time.Time `json:"createdAt"`
}

type Comment struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"postId"`
	ParentID  *int64    `json:"parentId,omitempty"`
	Author    string    `json:"author"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}
