package service

import (
	"context"

	"MISPRIS/internal/domain"
	"MISPRIS/internal/repository"
)

type HoPositionServiceImpl struct {
	repo repository.HoPositionRepository
}

func NewHoPositionService(repo repository.HoPositionRepository) *HoPositionServiceImpl {
	return &HoPositionServiceImpl{repo: repo}
}

func (s *HoPositionServiceImpl) ListByHo(ctx context.Context, hoID string) ([]*domain.HoPositionFull, error) {
	return s.repo.ListByHo(ctx, hoID)
}

// Create вызывает SQL-процедуру add_ho_position, которая автоматически пересчитывает итог.
func (s *HoPositionServiceImpl) Create(ctx context.Context, hoID, emobileID string, quantity int, unitPrice float64, note string) (string, error) {
	item := &domain.HoPosition{
		HoID:      hoID,
		EmobileID: emobileID,
		Quantity:  quantity,
		UnitPrice: unitPrice,
		Note:      note,
	}
	return s.repo.Create(ctx, item)
}

func (s *HoPositionServiceImpl) Update(ctx context.Context, id string, quantity int, unitPrice float64, note string) error {
	item := &domain.HoPosition{
		ID:        id,
		Quantity:  quantity,
		UnitPrice: unitPrice,
		Note:      note,
	}
	return s.repo.Update(ctx, item)
}

func (s *HoPositionServiceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
