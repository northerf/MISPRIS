package service

import (
	"context"

	"MISPRIS/internal/domain"
	"MISPRIS/internal/repository"
)

type HoClassRoleServiceImpl struct {
	repo repository.HoClassRoleRepository
}

func NewHoClassRoleService(repo repository.HoClassRoleRepository) *HoClassRoleServiceImpl {
	return &HoClassRoleServiceImpl{repo: repo}
}

func (s *HoClassRoleServiceImpl) List(ctx context.Context, hoClassID string) ([]*domain.HoClassRole, error) {
	return s.repo.List(ctx, hoClassID)
}

func (s *HoClassRoleServiceImpl) ListByClass(ctx context.Context, hoClassID string) ([]*domain.HoClassRole, error) {
	return s.repo.ListByClass(ctx, hoClassID)
}

func (s *HoClassRoleServiceImpl) Create(ctx context.Context, hoClassID, hoRoleID string, isRequired bool) (string, error) {
	item := &domain.HoClassRole{
		HoClassID:  hoClassID,
		HoRoleID:   hoRoleID,
		IsRequired: isRequired,
	}
	return s.repo.Create(ctx, item)
}

func (s *HoClassRoleServiceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
