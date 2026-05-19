package service

import (
	"context"
	"fmt"

	"MISPRIS/internal/domain"
	"MISPRIS/internal/repository"
)

type ShdServiceImpl struct {
	repo repository.ShdRepository
}

func NewShdService(repo repository.ShdRepository) *ShdServiceImpl {
	return &ShdServiceImpl{repo: repo}
}

func (s *ShdServiceImpl) List(ctx context.Context) ([]*domain.SHD, error) {
	return s.repo.List(ctx)
}

func (s *ShdServiceImpl) GetByID(ctx context.Context, id string) (*domain.SHD, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ShdServiceImpl) Create(ctx context.Context, name, shdType, inn, description string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("name is required")
	}
	item := &domain.SHD{
		Name:        name,
		ShdType:     shdType,
		INN:         inn,
		Description: description,
	}
	return s.repo.Create(ctx, item)
}

func (s *ShdServiceImpl) Update(ctx context.Context, id, name, shdType, inn, description string) error {
	item := &domain.SHD{
		ID:          id,
		Name:        name,
		ShdType:     shdType,
		INN:         inn,
		Description: description,
	}
	return s.repo.Update(ctx, item)
}

func (s *ShdServiceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
