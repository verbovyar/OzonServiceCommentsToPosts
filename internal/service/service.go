package service

import (
	"context"
	"errors"
	"ozonProject/internal/models"
	"ozonProject/internal/storage"
	"ozonProject/internal/utils"
	"ozonProject/internal/validation"
)

type Service struct {
	storage storage.Storage
}

func New(storage storage.Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) ListPosts(ctx context.Context, limit, offset *int) ([]*models.Post, error) {
	return s.storage.GetPosts(ctx, utils.ValueOrDefault(limit, 0), utils.ValueOrDefault(offset, 0))
}

func (s *Service) GetPost(ctx context.Context, id string) (*models.Post, error) {
	return s.storage.GetPostByID(ctx, id)
}

func (s *Service) CreatePost(ctx context.Context, title, content, author string, commentsEnabled *bool) (*models.Post, error) {
	return s.storage.CreatePost(ctx, title, content, author, utils.ValueOrDefault(commentsEnabled, false))
}

func (s *Service) CreateComment(ctx context.Context, postId string, parentId *string, author, content string) (*models.Comment, error) {
	if err := validation.ValidateCommentBody(content); err != nil {
		return nil, err
	}

	if err := s.storage.EnsureCommentsEnabled(ctx, postId); err != nil {
		return nil, validation.ErrCommentsOff
	}

	return s.storage.CreateComment(ctx, postId, utils.ValueOrDefault(parentId, ""), author, content)
}

func (s *Service) ListComments(ctx context.Context, postId string, parentId *string, limit, offset *int) ([]*models.Comment, error) {
	return s.storage.GetComments(ctx, postId,
		utils.ValueOrDefault(parentId, ""), utils.ValueOrDefault(limit, 0), utils.ValueOrDefault(offset, 0))
}

func ToUserError(err error) error {
	switch {
	case errors.Is(err, validation.ErrCommentsOff):
		return err
	case errors.Is(err, validation.ErrTooLong):
		return err
	case errors.Is(err, validation.ErrEmptyContent):
		return err
	default:
		return err
	}
}
