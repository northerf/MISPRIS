package service

import (
	"context"

	"MISPRIS/internal/domain"
	"MISPRIS/internal/repository"
)

type HoParameterValueServiceImpl struct {
	repo repository.HoParameterValueRepository
}

func NewHoParameterValueService(repo repository.HoParameterValueRepository) *HoParameterValueServiceImpl {
	return &HoParameterValueServiceImpl{repo: repo}
}

func (s *HoParameterValueServiceImpl) ListByHo(ctx context.Context, hoID string) ([]*domain.HoParameterValueFull, error) {
	return s.repo.ListByHo(ctx, hoID)
}

func (s *HoParameterValueServiceImpl) GetByID(ctx context.Context, id string) (*domain.HoParameterValue, error) {
	return s.repo.GetByID(ctx, id)
}

// Create вызывает SQL-процедуру write_ho_par для записи значения с валидацией.
func (s *HoParameterValueServiceImpl) Create(ctx context.Context, hoID, hoClassParameterID string, valReal float64, valInt int, valStr, valDate, enumValID string) (string, error) {
	item := &domain.HoParameterValue{
		HoID:               hoID,
		HoClassParameterID: hoClassParameterID,
		ValReal:            valReal,
		ValInt:             valInt,
		ValStr:             valStr,
		ValDate:            valDate,
		EnumValID:          enumValID,
	}
	return s.repo.Create(ctx, item)
}

// Update вызывает SQL-процедуру write_ho_par повторно (UPSERT-семантика).
func (s *HoParameterValueServiceImpl) Update(ctx context.Context, id, hoID, hoClassParameterID string, valReal float64, valInt int, valStr, valDate, enumValID string) error {
	item := &domain.HoParameterValue{
		ID:                 id,
		HoID:               hoID,
		HoClassParameterID: hoClassParameterID,
		ValReal:            valReal,
		ValInt:             valInt,
		ValStr:             valStr,
		ValDate:            valDate,
		EnumValID:          enumValID,
	}
	return s.repo.Update(ctx, item)
}

func (s *HoParameterValueServiceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
