package service

import (
	"context"
	"fmt"

	"MISPRIS/internal/domain"
	"MISPRIS/internal/repository"
)

type HoRoleServiceImpl struct {
	repo repository.HoRoleRepository
}

func NewHoRoleService(repo repository.HoRoleRepository) *HoRoleServiceImpl {
	return &HoRoleServiceImpl{repo: repo}
}

func (s *HoRoleServiceImpl) List(ctx context.Context) ([]*domain.HoRole, error) {
	return s.repo.List(ctx)
}

func (s *HoRoleServiceImpl) GetByID(ctx context.Context, id string) (*domain.HoRole, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *HoRoleServiceImpl) Create(ctx context.Context, name, description string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("name is required")
	}
	item := &domain.HoRole{
		Name:        name,
		Description: description,
	}
	return s.repo.Create(ctx, item)
}

func (s *HoRoleServiceImpl) Update(ctx context.Context, id, name, description string) error {
	item := &domain.HoRole{
		ID:          id,
		Name:        name,
		Description: description,
	}
	return s.repo.Update(ctx, item)
}

func (s *HoRoleServiceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
