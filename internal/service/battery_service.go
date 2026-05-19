package service

import (
	"context"
	"time"

	"MISPRIS/internal/domain"
	"MISPRIS/internal/repository"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type BatteryServiceImpl struct {
	repo repository.BatteryRepository
	db   *sqlx.DB
}

func NewBatteryService(repo repository.BatteryRepository, db *sqlx.DB) *BatteryServiceImpl {
	return &BatteryServiceImpl{repo: repo, db: db}
}

func (s *BatteryServiceImpl) Create(ctx context.Context, name string, batteryType string, batteryCapacity string, batteryInfo string) (string, error) {
	return s.repo.Create(ctx, &domain.Battery{
		ID:              uuid.New().String(),
		BatteryName:     name,
		BatteryInfo:     batteryInfo,
		BatteryCapacity: batteryCapacity,
		BatteryType:     batteryType,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	})
}

func (s *BatteryServiceImpl) Update(ctx context.Context, id string, name string, batteryType string, batteryCapacity string, batteryInfo string) error {
	return s.repo.Update(ctx, &domain.Battery{ID: id, BatteryName: name, BatteryType: batteryType, BatteryCapacity: batteryCapacity, BatteryInfo: batteryInfo, UpdatedAt: time.Now()})
}

func (s *BatteryServiceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *BatteryServiceImpl) List(ctx context.Context) ([]domain.Battery, error) {
	batteries, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Battery, len(batteries))
	for i, b := range batteries {
		result[i] = *b
	}
	return result, nil
}

func (s *BatteryServiceImpl) GetByID(ctx context.Context, id string) (*domain.Battery, error) {
	return s.repo.GetByID(ctx, id)
}
