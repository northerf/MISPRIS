package service

import (
	"context"
	"fmt"
	"strconv"

	"MISPRIS/internal/domain"
	"MISPRIS/internal/repository"
)

type ComponentParameterServiceImpl struct {
	repo repository.ComponentParameterRepository
}

func NewComponentParameterService(repo repository.ComponentParameterRepository) *ComponentParameterServiceImpl {
	return &ComponentParameterServiceImpl{repo: repo}
}

func (s *ComponentParameterServiceImpl) List(ctx context.Context) ([]*domain.ComponentParameter, error) {
	return s.repo.List(ctx)
}

func (s *ComponentParameterServiceImpl) GetByID(ctx context.Context, id string) (*domain.ComponentParameter, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ComponentParameterServiceImpl) Create(ctx context.Context,
	componentType, parameterID, orderNumStr, minValStr, maxValStr string) (string, error) {
	if componentType == "" || parameterID == "" {
		return "", fmt.Errorf("component_type and parameter_id are required")
	}
	orderNum, _ := strconv.Atoi(orderNumStr)
	minVal, _ := strconv.ParseFloat(minValStr, 64)
	maxVal, _ := strconv.ParseFloat(maxValStr, 64)
	cp := &domain.ComponentParameter{
		ComponentType: componentType,
		ParameterID:   parameterID,
		OrderNum:      orderNum,
		MinVal:        minVal,
		MaxVal:        maxVal,
	}
	return s.repo.Create(ctx, cp)
}

func (s *ComponentParameterServiceImpl) Update(ctx context.Context,
	id, orderNumStr, minValStr, maxValStr string) error {
	orderNum, _ := strconv.Atoi(orderNumStr)
	minVal, _ := strconv.ParseFloat(minValStr, 64)
	maxVal, _ := strconv.ParseFloat(maxValStr, 64)
	cp := &domain.ComponentParameter{
		ID:       id,
		OrderNum: orderNum,
		MinVal:   minVal,
		MaxVal:   maxVal,
	}
	return s.repo.Update(ctx, cp)
}

func (s *ComponentParameterServiceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *ComponentParameterServiceImpl) GetByType(ctx context.Context, componentType string) ([]*domain.ComponentParameterFull, error) {
	if componentType == "" {
		return nil, fmt.Errorf("component_type is required")
	}
	return s.repo.GetByType(ctx, componentType)
}

func (s *ComponentParameterServiceImpl) CopyFromType(ctx context.Context, fromType, toType string) error {
	if fromType == "" || toType == "" {
		return fmt.Errorf("from_type and to_type are required")
	}
	return s.repo.CopyFromType(ctx, fromType, toType)
}
