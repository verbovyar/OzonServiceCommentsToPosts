package validation

import "errors"

const MaxCommentLen = 2000

var (
	ErrTooLong      = errors.New("content too long")
	ErrCommentsOff  = errors.New("comments are disabled for this post")
	ErrEmptyContent = errors.New("content is empty")
)

func ValidateCommentBody(s string) error {
	if len(s) == 0 {
		return ErrEmptyContent
	}

	if len(s) > MaxCommentLen {
		return ErrTooLong
	}

	return nil
}
