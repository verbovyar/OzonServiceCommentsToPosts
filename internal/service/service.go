package service

import (
	"context"
	"errors"
	"ozonProject/internal/models"
	"ozonProject/internal/storage"
	"ozonProject/internal/validation"
	"strconv"
)

type Service struct {
	storage storage.Storage
}

func New(storage storage.Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) ListPosts(ctx context.Context, limit, offset int) ([]*models.Post, error) {
	if limit <= 0 {
		limit = 10
	}

	if offset < 0 {
		offset = 0
	}

	return s.storage.GetPosts(ctx, limit, offset)
}

func (s *Service) GetPost(ctx context.Context, idStr string) (*models.Post, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, err
	}

	return s.storage.GetPostByID(ctx, id)
}

func (s *Service) CreatePost(ctx context.Context, title, content, author string, commentsEnabled bool) (*models.Post, error) {
	return s.storage.CreatePost(ctx, title, content, author, commentsEnabled)
}

func (s *Service) CreateComment(ctx context.Context, postIDStr string, parentIDStr *string, author, content string) (*models.Comment, error) {
	if err := validation.ValidateCommentBody(content); err != nil {
		return nil, err
	}

	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		return nil, err
	}

	if err := s.storage.EnsureCommentsEnabled(ctx, postID); err != nil {
		return nil, validation.ErrCommentsOff
	}

	var parentID *int64
	if parentIDStr != nil {
		pid, err := strconv.ParseInt(*parentIDStr, 10, 64)
		if err != nil {
			return nil, err
		}
		parentID = &pid
	}

	return s.storage.CreateComment(ctx, postID, parentID, author, content)
}

func (s *Service) ListComments(ctx context.Context, postIDStr string, parentIDStr *string, limit, offset int) ([]*models.Comment, error) {
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		return nil, err
	}

	var parentID *int64
	if parentIDStr != nil && *parentIDStr != "" {
		pid, err := strconv.ParseInt(*parentIDStr, 10, 64)
		if err != nil {
			return nil, err
		}
		parentID = &pid
	}

	if limit <= 0 {
		limit = 10
	}

	if offset < 0 {
		offset = 0
	}

	return s.storage.GetComments(ctx, postID, parentID, limit, offset)
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
