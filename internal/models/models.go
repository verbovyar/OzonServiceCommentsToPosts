package models

import "time"

type Post struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Content         string    `json:"content"`
	Author          string    `json:"author"`
	CommentsEnabled bool      `json:"commentsEnabled"`
	CreatedAt       time.Time `json:"createdAt"`
}

type Comment struct {
	ID        string    `json:"id"`
	PostID    string    `json:"postId"`
	ParentID  *string   `json:"parentId,omitempty"`
	Author    string    `json:"author"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}
