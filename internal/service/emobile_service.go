package service

import (
	"context"
	"fmt"
	"time"

	"MISPRIS/internal/domain"
	"MISPRIS/internal/repository"

	"github.com/jmoiron/sqlx"
)

type EmobileServiceImpl struct {
	db                   *sqlx.DB
	emobileRepo          repository.EmobileRepository
	chargerSystemService ChargerSystemService
	bodyService          BodyService
	electronicsService   ElectronicsService
	chassisService       ChassisService
	batteryService       BatteryService
	powerPointService    PowerPointService
}

func NewEmobileService(
	db *sqlx.DB,
	emobileRepo repository.EmobileRepository,
	chargerSystemService ChargerSystemService,
	bodyService BodyService,
	electronicsService ElectronicsService,
	chassisService ChassisService,
	batteryService BatteryService,
	powerPointService PowerPointService,
) *EmobileServiceImpl {
	return &EmobileServiceImpl{
		db:                   db,
		emobileRepo:          emobileRepo,
		chargerSystemService: chargerSystemService,
		bodyService:          bodyService,
		electronicsService:   electronicsService,
		chassisService:       chassisService,
		batteryService:       batteryService,
		powerPointService:    powerPointService,
	}
}

func (s *EmobileServiceImpl) Create(ctx context.Context, name string, powerPointID string,
	batteryID string, chargerSystemID string, chassisID string, bodyID string, electonicsID string) (string, error) {

	if _, err := s.chargerSystemService.GetByID(ctx, chargerSystemID); err != nil {
		return "", fmt.Errorf("chargerSystem %s not found: %w", chargerSystemID, err)
	}

	if _, err := s.bodyService.GetByID(ctx, bodyID); err != nil {
		return "", fmt.Errorf("body %s not found: %w", bodyID, err)
	}

	if _, err := s.electronicsService.GetByID(ctx, electonicsID); err != nil {
		return "", fmt.Errorf("electonics %s not found: %w", electonicsID, err)
	}

	if _, err := s.chassisService.GetByID(ctx, chassisID); err != nil {
		return "", fmt.Errorf("chassis %s not found: %w", chassisID, err)
	}

	if _, err := s.batteryService.GetByID(ctx, batteryID); err != nil {
		return "", fmt.Errorf("battery %s not found: %w", batteryID, err)
	}

	if _, err := s.powerPointService.GetByID(ctx, powerPointID); err != nil {
		return "", fmt.Errorf("powerPoint %s not found: %w", powerPointID, err)
	}
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	id, err := s.emobileRepo.Create(ctx, tx, &domain.Emobile{
		Name:            name,
		ChargerSystemID: chargerSystemID,
		BodyID:          bodyID,
		ElectronicsID:   electonicsID,
		ChassisID:       chassisID,
		BatteryID:       batteryID,
		PowerPointID:    powerPointID,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	})
	if err != nil {
		return "", err
	}
	return id, tx.Commit()
}

func (s *EmobileServiceImpl) Update(ctx context.Context, id string, name string, powerPointID string,
	batteryID string, chargerSystemID string, chassisID string, bodyID string, electonicsID string) error {
	return s.emobileRepo.Update(ctx, &domain.Emobile{
		ID:              id,
		Name:            name,
		ChargerSystemID: chargerSystemID,
		BodyID:          bodyID,
		ElectronicsID:   electonicsID,
		ChassisID:       chassisID,
		BatteryID:       batteryID,
		PowerPointID:    powerPointID,
		UpdatedAt:       time.Now(),
	})
}

func (s *EmobileServiceImpl) Delete(ctx context.Context, id string) error {
	return s.emobileRepo.Delete(ctx, id)
}

func (s *EmobileServiceImpl) List(ctx context.Context) ([]domain.Emobile, error) {
	rows, err := s.emobileRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Emobile, len(rows))
	for i, r := range rows {
		result[i] = *r
	}
	return result, nil
}

func (s *EmobileServiceImpl) GetByID(ctx context.Context, id string) (*domain.Emobile, error) {
	return s.emobileRepo.GetByID(ctx, id)
}
