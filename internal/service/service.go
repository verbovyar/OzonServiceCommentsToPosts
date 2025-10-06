package service

import (
	"context"
	"ozonProject/internal/models"
)

type Service struct {
	// store
}

func New() *Service {
	return &Service{}
}

func (s *Service) ListPosts(ctx context.Context, limit, offset int) ([]*models.Post, error) {
	return nil, nil
}

func (s *Service) GetPost(ctx context.Context, idStr string) (*models.Post, error) {
	return nil, nil
}

func (s *Service) CreatePost(ctx context.Context, title, content, author string, commentsEnabled bool) (*models.Post, error) {
	return nil, nil
}

func (s *Service) CreateComment(ctx context.Context, postIDStr string, parentIDStr *string, author, content string) (*models.Comment, error) {
	return nil, nil
}

func (s *Service) ListComments(ctx context.Context, postIDStr string, parentIDStr *string, limit, offset int) ([]*models.Comment, error) {
	return nil, nil
}
