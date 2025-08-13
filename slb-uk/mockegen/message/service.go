package message

import (
	"context"
	"errors"
	"strings"
)

var (
    ErrEmptyContent = errors.New("content cannot be empty")
    ErrInvalidID    = errors.New("id must be > 0")
)

type Service struct {
    repo Repository
}

func NewService(r Repository) *Service { return &Service{repo: r} }

func (s *Service) Create(ctx context.Context, content string) (Message, error) {
    content = strings.TrimSpace(content)
    if content == "" {
        return Message{}, ErrEmptyContent
    }
    return s.repo.Create(ctx, Message{Content: content})
}

func (s *Service) Get(ctx context.Context, id int) (Message, error) {
    if id <= 0 {
        return Message{}, ErrInvalidID
    }
    return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int, content string) (Message, error) {
    if id <= 0 {
        return Message{}, ErrInvalidID
    }
    content = strings.TrimSpace(content)
    if content == "" {
        return Message{}, ErrEmptyContent
    }
    return s.repo.Update(ctx, Message{ID: id, Content: content})
}

func (s *Service) Delete(ctx context.Context, id int) error {
    if id <= 0 {
        return ErrInvalidID
    }
    return s.repo.Delete(ctx, id)
}
