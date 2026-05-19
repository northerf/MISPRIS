package service

import (
	"context"
	"fmt"

	"MISPRIS/internal/domain"
	"MISPRIS/internal/repository"
)

type DocumentClassServiceImpl struct {
	repo repository.DocumentClassRepository
}

func NewDocumentClassService(repo repository.DocumentClassRepository) *DocumentClassServiceImpl {
	return &DocumentClassServiceImpl{repo: repo}
}

func (s *DocumentClassServiceImpl) List(ctx context.Context) ([]*domain.DocumentClass, error) {
	return s.repo.List(ctx)
}

func (s *DocumentClassServiceImpl) GetByID(ctx context.Context, id string) (*domain.DocumentClass, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *DocumentClassServiceImpl) Create(ctx context.Context, name, code, description string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("name is required")
	}
	item := &domain.DocumentClass{
		Name:        name,
		Code:        code,
		Description: description,
	}
	return s.repo.Create(ctx, item)
}

func (s *DocumentClassServiceImpl) Update(ctx context.Context, id, name, code, description string) error {
	item := &domain.DocumentClass{
		ID:          id,
		Name:        name,
		Code:        code,
		Description: description,
	}
	return s.repo.Update(ctx, item)
}

func (s *DocumentClassServiceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
