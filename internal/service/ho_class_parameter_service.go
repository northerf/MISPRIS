package service

import (
	"context"

	"MISPRIS/internal/domain"
	"MISPRIS/internal/repository"
)

type HoClassParameterServiceImpl struct {
	repo repository.HoClassParameterRepository
}

func NewHoClassParameterService(repo repository.HoClassParameterRepository) *HoClassParameterServiceImpl {
	return &HoClassParameterServiceImpl{repo: repo}
}

func (s *HoClassParameterServiceImpl) List(ctx context.Context, hoClassID string) ([]*domain.HoClassParameter, error) {
	return s.repo.List(ctx, hoClassID)
}

func (s *HoClassParameterServiceImpl) GetByID(ctx context.Context, id string) (*domain.HoClassParameter, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *HoClassParameterServiceImpl) Create(ctx context.Context, hoClassID, parameterID string, orderNum int, minVal, maxVal float64) (string, error) {
	item := &domain.HoClassParameter{
		HoClassID:   hoClassID,
		ParameterID: parameterID,
		OrderNum:    orderNum,
		MinVal:      minVal,
		MaxVal:      maxVal,
	}
	return s.repo.Create(ctx, item)
}

func (s *HoClassParameterServiceImpl) Update(ctx context.Context, id string, orderNum int, minVal, maxVal float64) error {
	item := &domain.HoClassParameter{
		ID:       id,
		OrderNum: orderNum,
		MinVal:   minVal,
		MaxVal:   maxVal,
	}
	return s.repo.Update(ctx, item)
}

func (s *HoClassParameterServiceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// GetByHoClass вызывает SQL-функцию get_ho_class_parameters.
func (s *HoClassParameterServiceImpl) GetByHoClass(ctx context.Context, hoClassID string) ([]*domain.HoClassParameterFull, error) {
	return s.repo.GetByHoClass(ctx, hoClassID)
}

// CopyFromClass копирует параметры из одного класса ХО в другой.
func (s *HoClassParameterServiceImpl) CopyFromClass(ctx context.Context, fromClassID, toClassID string) error {
	return s.repo.CopyFromClass(ctx, fromClassID, toClassID)
}
