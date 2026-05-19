package service

import (
	"context"
	"fmt"

	"MISPRIS/internal/domain"
	"MISPRIS/internal/repository"
)

type HoClassServiceImpl struct {
	repo repository.HoClassRepository
}

func NewHoClassService(repo repository.HoClassRepository) *HoClassServiceImpl {
	return &HoClassServiceImpl{repo: repo}
}

func (s *HoClassServiceImpl) List(ctx context.Context) ([]*domain.HoClass, error) {
	return s.repo.List(ctx)
}

func (s *HoClassServiceImpl) GetByID(ctx context.Context, id string) (*domain.HoClass, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *HoClassServiceImpl) Create(ctx context.Context, name, designation, parentID string, isTerminal bool) (string, error) {
	if name == "" {
		return "", fmt.Errorf("name is required")
	}
	item := &domain.HoClass{
		Name:        name,
		Designation: designation,
		ParentID:    parentID,
		IsTerminal:  isTerminal,
	}
	return s.repo.Create(ctx, item)
}

func (s *HoClassServiceImpl) Update(ctx context.Context, id, name, designation, parentID string, isTerminal bool) error {
	item := &domain.HoClass{
		ID:          id,
		Name:        name,
		Designation: designation,
		ParentID:    parentID,
		IsTerminal:  isTerminal,
	}
	return s.repo.Update(ctx, item)
}

func (s *HoClassServiceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *HoClassServiceImpl) GetTerminal(ctx context.Context) ([]*domain.HoClass, error) {
	return s.repo.GetTerminal(ctx)
}

func (s *HoClassServiceImpl) GetChildren(ctx context.Context, parentID string) ([]*domain.HoClass, error) {
	return s.repo.GetChildren(ctx, parentID)
}
