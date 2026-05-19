package service

import (
	"context"
	"fmt"
	"strconv"

	"MISPRIS/internal/domain"
	"MISPRIS/internal/repository"
)

type EnumClassServiceImpl struct {
	repo repository.EnumClassRepository
}

func NewEnumClassService(repo repository.EnumClassRepository) *EnumClassServiceImpl {
	return &EnumClassServiceImpl{repo: repo}
}

func (s *EnumClassServiceImpl) List(ctx context.Context) ([]*domain.EnumClass, error) {
	return s.repo.List(ctx)
}

func (s *EnumClassServiceImpl) GetByID(ctx context.Context, id string) (*domain.EnumClass, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *EnumClassServiceImpl) Create(ctx context.Context, name, componentType string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("name is required")
	}
	ec := &domain.EnumClass{
		Name:          name,
		ComponentType: componentType,
	}
	return s.repo.Create(ctx, ec)
}

func (s *EnumClassServiceImpl) Update(ctx context.Context, id, name, componentType string) error {
	ec := &domain.EnumClass{
		ID:            id,
		Name:          name,
		ComponentType: componentType,
	}
	return s.repo.Update(ctx, ec)
}

func (s *EnumClassServiceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *EnumClassServiceImpl) GetValues(ctx context.Context, id string) ([]*domain.EnumPosition, error) {
	return s.repo.GetValues(ctx, id)
}

func (s *EnumClassServiceImpl) ValidateValue(ctx context.Context, enumClassID, value string) (bool, error) {
	if _, err := strconv.Atoi(enumClassID); err != nil {
		return false, fmt.Errorf("invalid enum_class_id: %s", enumClassID)
	}
	return s.repo.ValidateValue(ctx, enumClassID, value)
}
