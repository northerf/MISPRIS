package service

import (
	"context"
	"fmt"

	"MISPRIS/internal/domain"
	"MISPRIS/internal/repository"
)

type ParameterServiceImpl struct {
	repo repository.ParameterRepository
}

func NewParameterService(repo repository.ParameterRepository) *ParameterServiceImpl {
	return &ParameterServiceImpl{repo: repo}
}

func (s *ParameterServiceImpl) List(ctx context.Context) ([]*domain.Parameter, error) {
	return s.repo.List(ctx)
}

func (s *ParameterServiceImpl) GetByID(ctx context.Context, id string) (*domain.Parameter, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ParameterServiceImpl) Create(ctx context.Context,
	designation, name, paramType, measuringUnit, enumClassID string) (string, error) {
	if designation == "" || name == "" || paramType == "" {
		return "", fmt.Errorf("designation, name and param_type are required")
	}
	validTypes := map[string]bool{"real": true, "int": true, "str": true, "enum": true}
	if !validTypes[paramType] {
		return "", fmt.Errorf("param_type must be one of: real, int, str, enum")
	}
	p := &domain.Parameter{
		Designation:   designation,
		Name:          name,
		ParamType:     paramType,
		MeasuringUnit: measuringUnit,
		EnumClassID:   enumClassID,
	}
	return s.repo.Create(ctx, p)
}

func (s *ParameterServiceImpl) Update(ctx context.Context,
	id, designation, name, paramType, measuringUnit, enumClassID string) error {
	p := &domain.Parameter{
		ID:            id,
		Designation:   designation,
		Name:          name,
		ParamType:     paramType,
		MeasuringUnit: measuringUnit,
		EnumClassID:   enumClassID,
	}
	return s.repo.Update(ctx, p)
}

func (s *ParameterServiceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
