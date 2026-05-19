package service

import (
	"MISPRIS/internal/domain"
	"MISPRIS/internal/repository"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

//Controllers

type ControllerServiceImpl struct {
	repo repository.ControllerRepository
}

func NewControllerService(repo repository.ControllerRepository) *ControllerServiceImpl {
	return &ControllerServiceImpl{repo: repo}
}

func (s *ControllerServiceImpl) Create(ctx context.Context, controllerName, controllerInfo string) (string, error) {
	return s.repo.Create(ctx, &domain.Controller{
		ID:        uuid.New().String(),
		Name:      controllerName,
		Info:      controllerInfo,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
}

func (s *ControllerServiceImpl) Update(ctx context.Context, id, controllerName, controllerInfo string) error {
	return s.repo.Update(ctx, &domain.Controller{ID: id, Name: controllerName, Info: controllerInfo, UpdatedAt: time.Now()})
}

func (s *ControllerServiceImpl) Delete(ctx context.Context, controllerID string) error {
	return s.repo.Delete(ctx, controllerID)
}

func (s *ControllerServiceImpl) List(ctx context.Context) ([]domain.Controller, error) {
	rows, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Controller, len(rows))
	for i, r := range rows {
		result[i] = *r
	}
	return result, nil
}

func (s *ControllerServiceImpl) GetByID(ctx context.Context, id string) (*domain.Controller, error) {
	return s.repo.GetByID(ctx, id)
}

//Sensors

type SensorServiceImpl struct {
	repo repository.SensorRepository
}

func NewSensorService(repo repository.SensorRepository) *SensorServiceImpl {
	return &SensorServiceImpl{repo: repo}
}

func (s *SensorServiceImpl) Create(ctx context.Context, sensorName, sensorInfo string) (string, error) {
	return s.repo.Create(ctx, &domain.Sensor{
		ID:        uuid.New().String(),
		Name:      sensorName,
		Info:      sensorInfo,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
}

func (s *SensorServiceImpl) Update(ctx context.Context, sensorID, sensorName, sensorInfo string) error {
	return s.repo.Update(ctx, &domain.Sensor{
		ID:        sensorID,
		Name:      sensorName,
		Info:      sensorInfo,
		UpdatedAt: time.Now(),
	})
}

func (s *SensorServiceImpl) Delete(ctx context.Context, sensorID string) error {
	return s.repo.Delete(ctx, sensorID)
}

func (s *SensorServiceImpl) List(ctx context.Context) ([]domain.Sensor, error) {
	rows, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Sensor, len(rows))
	for i, r := range rows {
		result[i] = *r
	}
	return result, nil
}

func (s *SensorServiceImpl) GetByID(ctx context.Context, sensorID string) (*domain.Sensor, error) {
	return s.repo.GetByID(ctx, sensorID)
}

//Wiring

type WiringServiceImpl struct {
	repo repository.WiringRepository
}

func NewWiringService(repo repository.WiringRepository) *WiringServiceImpl {
	return &WiringServiceImpl{repo: repo}
}

func (s *WiringServiceImpl) Create(ctx context.Context, wiringName, wiringInfo string) (string, error) {
	return s.repo.Create(ctx, &domain.Wiring{
		ID:        uuid.New().String(),
		Name:      wiringName,
		Info:      wiringInfo,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
}

func (s *WiringServiceImpl) Update(ctx context.Context, wiringID, wiringName, wiringInfo string) error {
	return s.repo.Update(ctx, &domain.Wiring{
		ID:        wiringID,
		Name:      wiringName,
		Info:      wiringInfo,
		UpdatedAt: time.Now(),
	})
}

func (s *WiringServiceImpl) Delete(ctx context.Context, sensorID string) error {
	return s.repo.Delete(ctx, sensorID)
}

func (s *WiringServiceImpl) List(ctx context.Context) ([]domain.Wiring, error) {
	rows, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Wiring, len(rows))
	for i, r := range rows {
		result[i] = *r
	}
	return result, nil
}

func (s *WiringServiceImpl) GetByID(ctx context.Context, sensorID string) (*domain.Wiring, error) {
	return s.repo.GetByID(ctx, sensorID)
}

//Electronics

type ElectronicsServiceImpl struct {
	db                *sqlx.DB
	electronicsRepo   repository.ElectronicsRepository
	wiringService     WiringService
	controllerService ControllerService
	sensorService     SensorService
}

func NewElectronicsService(
	db *sqlx.DB,
	electronicsRepo repository.ElectronicsRepository,
	wiringService WiringService,
	controllerService ControllerService,
	sensorService SensorService,
) *ElectronicsServiceImpl {
	return &ElectronicsServiceImpl{
		db:                db,
		electronicsRepo:   electronicsRepo,
		wiringService:     wiringService,
		controllerService: controllerService,
		sensorService:     sensorService,
	}
}

func (s *ElectronicsServiceImpl) Create(ctx context.Context, wiringID, sensorID, controllerID string) (string, error) {
	if _, err := s.controllerService.GetByID(ctx, controllerID); err != nil {
		return "", fmt.Errorf("carcass %s not found: %w", controllerID, err)
	}
	if _, err := s.sensorService.GetByID(ctx, sensorID); err != nil {
		return "", fmt.Errorf("doors %s not found: %w", sensorID, err)
	}
	if _, err := s.wiringService.GetByID(ctx, wiringID); err != nil {
		return "", fmt.Errorf("wings %s not found: %w", wiringID, err)
	}
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	id, err := s.electronicsRepo.CreateTx(ctx, tx, &domain.Electronics{
		ID:           uuid.New().String(),
		WiringID:     wiringID,
		SensorID:     sensorID,
		ControllerID: controllerID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	})
	if err != nil {
		return "", err
	}
	return id, tx.Commit()
}

func (s *ElectronicsServiceImpl) Update(ctx context.Context, id, wiringID, sensorID, controllerID string) error {
	return s.electronicsRepo.Update(ctx, &domain.Electronics{
		ID:           id,
		WiringID:     wiringID,
		SensorID:     sensorID,
		ControllerID: controllerID,
		UpdatedAt:    time.Now(),
	})
}

func (s *ElectronicsServiceImpl) Delete(ctx context.Context, id string) error {
	return s.electronicsRepo.Delete(ctx, id)
}

func (s *ElectronicsServiceImpl) List(ctx context.Context) ([]domain.Electronics, error) {
	rows, err := s.electronicsRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Electronics, len(rows))
	for i, r := range rows {
		result[i] = *r
	}
	return result, nil
}

func (s *ElectronicsServiceImpl) GetByID(ctx context.Context, id string) (*domain.Electronics, error) {
	return s.electronicsRepo.GetByID(ctx, id)
}
