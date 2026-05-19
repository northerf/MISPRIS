package service

import (
	"context"
	"fmt"

	"MISPRIS/internal/domain"
	"MISPRIS/internal/repository"
)

type HoInstanceServiceImpl struct {
	repo repository.HoInstanceRepository
}

func NewHoInstanceService(repo repository.HoInstanceRepository) *HoInstanceServiceImpl {
	return &HoInstanceServiceImpl{repo: repo}
}

func (s *HoInstanceServiceImpl) List(ctx context.Context, hoClassID string) ([]*domain.HoInstance, error) {
	return s.repo.List(ctx, hoClassID)
}

func (s *HoInstanceServiceImpl) GetByID(ctx context.Context, id string) (*domain.HoInstance, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *HoInstanceServiceImpl) Create(ctx context.Context, hoClassID, docNumber, docDate string, totalAmount float64, note string) (string, error) {
	if hoClassID == "" {
		return "", fmt.Errorf("ho_class_id is required")
	}
	item := &domain.HoInstance{
		HoClassID:   hoClassID,
		DocNumber:   docNumber,
		DocDate:     docDate,
		TotalAmount: totalAmount,
		Note:        note,
	}
	return s.repo.Create(ctx, item)
}

func (s *HoInstanceServiceImpl) Update(ctx context.Context, id, status, docNumber, docDate string, totalAmount float64, note string) error {
	item := &domain.HoInstance{
		ID:          id,
		Status:      status,
		DocNumber:   docNumber,
		DocDate:     docDate,
		TotalAmount: totalAmount,
		Note:        note,
	}
	return s.repo.Update(ctx, item)
}

func (s *HoInstanceServiceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// FindByClass вызывает SQL-функцию find_ho_by_class.
func (s *HoInstanceServiceImpl) FindByClass(ctx context.Context, hoClassID string) ([]*domain.HoInstanceFull, error) {
	return s.repo.FindByClass(ctx, hoClassID)
}
