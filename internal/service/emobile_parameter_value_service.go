package service

import (
	"context"
	"fmt"
	"strconv"

	"MISPRIS/internal/domain"
	"MISPRIS/internal/repository"
)

type EmobileParameterValueServiceImpl struct {
	repo repository.EmobileParameterValueRepository
}

func NewEmobileParameterValueService(repo repository.EmobileParameterValueRepository) *EmobileParameterValueServiceImpl {
	return &EmobileParameterValueServiceImpl{repo: repo}
}

func (s *EmobileParameterValueServiceImpl) List(ctx context.Context) ([]*domain.EmobileParameterValue, error) {
	return s.repo.List(ctx)
}

func (s *EmobileParameterValueServiceImpl) GetByID(ctx context.Context, id string) (*domain.EmobileParameterValue, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *EmobileParameterValueServiceImpl) Create(ctx context.Context,
	emobileID, componentParameterID, valRealStr, valIntStr, valStr, enumValID string) (string, error) {
	if emobileID == "" || componentParameterID == "" {
		return "", fmt.Errorf("emobile_id and component_parameter_id are required")
	}
	valReal, _ := strconv.ParseFloat(valRealStr, 64)
	valInt, _ := strconv.Atoi(valIntStr)
	v := &domain.EmobileParameterValue{
		EmobileID:            emobileID,
		ComponentParameterID: componentParameterID,
		ValReal:              valReal,
		ValInt:               valInt,
		ValStr:               valStr,
		EnumValID:            enumValID,
	}
	return s.repo.Create(ctx, v)
}

func (s *EmobileParameterValueServiceImpl) Update(ctx context.Context,
	id, valRealStr, valIntStr, valStr, enumValID string) error {
	valReal, _ := strconv.ParseFloat(valRealStr, 64)
	valInt, _ := strconv.Atoi(valIntStr)
	v := &domain.EmobileParameterValue{
		ID:        id,
		ValReal:   valReal,
		ValInt:    valInt,
		ValStr:    valStr,
		EnumValID: enumValID,
	}
	return s.repo.Update(ctx, v)
}

func (s *EmobileParameterValueServiceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *EmobileParameterValueServiceImpl) GetByEmobile(ctx context.Context, emobileID string) ([]*domain.EmobileParameterValue, error) {
	if emobileID == "" {
		return nil, fmt.Errorf("emobile_id is required")
	}
	return s.repo.GetByEmobile(ctx, emobileID)
}
