package service

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"MISPRIS/internal/domain"
	"MISPRIS/internal/repository"

	"github.com/google/uuid"
)

// Power Point
type PowerPointServiceImpl struct {
	repo            repository.PowerPointRepository
	db              *sqlx.DB
	engineService   EngineService
	inverterService InverterService
	gearBoxService  GearboxService
}

func NewPowerPointService(repo repository.PowerPointRepository, db *sqlx.DB,

	engineService EngineService,
	inverterService InverterService,
	gearBoxService GearboxService,
) *PowerPointServiceImpl {
	return &PowerPointServiceImpl{
		repo:            repo,
		db:              db,
		engineService:   engineService,
		inverterService: inverterService,
		gearBoxService:  gearBoxService,
	}
}

func (s *PowerPointServiceImpl) Create(ctx context.Context, engineID, invertorID, gearboxID string) (string, error) {
	if _, err := s.engineService.GetByID(ctx, engineID); err != nil {
		return "", fmt.Errorf("engine %s not found: %w", engineID, err)
	}
	if _, err := s.inverterService.GetByID(ctx, invertorID); err != nil {
		return "", fmt.Errorf("inverter %s not found: %w", invertorID, err)
	}
	if _, err := s.gearBoxService.GetByID(ctx, gearboxID); err != nil {
		return "", fmt.Errorf("wings %s not found: %w", gearboxID, err)
	}
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	id, err := s.repo.CreateTx(ctx, tx, &domain.PowerPoint{
		ID:         uuid.New().String(),
		EngineID:   engineID,
		InverterID: invertorID,
		GearboxID:  gearboxID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	})
	if err != nil {
		return "", err
	}
	return id, tx.Commit()
}

func (s *PowerPointServiceImpl) Update(ctx context.Context, id string, engineID, invertorID, gearboxID string) error {
	return s.repo.Update(ctx, &domain.PowerPoint{ID: id, EngineID: engineID, InverterID: invertorID, GearboxID: gearboxID, UpdatedAt: time.Now()})
}

func (s *PowerPointServiceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *PowerPointServiceImpl) List(ctx context.Context) ([]domain.PowerPoint, error) {
	rows, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.PowerPoint, len(rows))
	for i, r := range rows {
		result[i] = *r
	}
	return result, nil
}

func (s *PowerPointServiceImpl) GetByID(ctx context.Context, id string) (*domain.PowerPoint, error) {
	return s.repo.GetByID(ctx, id)
}

// Engine
type EngineServiceImpl struct {
	repo repository.EngineRepository
}

func NewEngineService(repo repository.EngineRepository) *EngineServiceImpl {
	return &EngineServiceImpl{repo: repo}
}

func (s *EngineServiceImpl) Create(ctx context.Context, name, engineType, info string) (string, error) {
	e := &domain.Engine{
		ID:         uuid.New().String(),
		EngineName: name,
		EngineType: engineType,
		EngineInfo: info,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	return s.repo.Create(ctx, e)
}

func (s *EngineServiceImpl) Update(ctx context.Context, id, name, engineType, info string) error {
	return s.repo.Update(ctx, &domain.Engine{ID: id, EngineName: name, EngineType: engineType, EngineInfo: info, UpdatedAt: time.Now()})
}

func (s *EngineServiceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *EngineServiceImpl) List(ctx context.Context) ([]domain.Engine, error) {
	engines, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Engine, len(engines))
	for i, e := range engines {
		result[i] = *e
	}
	return result, nil
}

func (s *EngineServiceImpl) GetByID(ctx context.Context, id string) (*domain.Engine, error) {
	return s.repo.GetByID(ctx, id)
}

// Inverter
type InverterServiceImpl struct {
	repo repository.InverterRepository
}

func NewInverterService(repo repository.InverterRepository) *InverterServiceImpl {
	return &InverterServiceImpl{repo: repo}
}

func (s *InverterServiceImpl) Create(ctx context.Context, name, info string) (string, error) {
	i := &domain.Inverter{
		ID:           uuid.New().String(),
		InverterName: name,
		InverterInfo: info,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	return s.repo.Create(ctx, i)
}

func (s *InverterServiceImpl) Update(ctx context.Context, id, name, info string) error {
	return s.repo.Update(ctx, &domain.Inverter{ID: id, InverterName: name, InverterInfo: info, UpdatedAt: time.Now()})
}

func (s *InverterServiceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *InverterServiceImpl) List(ctx context.Context) ([]domain.Inverter, error) {
	inverters, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Inverter, len(inverters))
	for i, inv := range inverters {
		result[i] = *inv
	}
	return result, nil
}

func (s *InverterServiceImpl) GetByID(ctx context.Context, id string) (*domain.Inverter, error) {
	return s.repo.GetByID(ctx, id)
}

// Gearbox
type GearboxServiceImpl struct {
	repo repository.GearboxRepository
}

func NewGearboxService(repo repository.GearboxRepository) *GearboxServiceImpl {
	return &GearboxServiceImpl{repo: repo}
}

func (s *GearboxServiceImpl) Create(ctx context.Context, name, info string) (string, error) {
	g := &domain.Gearbox{
		ID:          uuid.New().String(),
		GearboxName: name,
		GearboxInfo: info,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return s.repo.Create(ctx, g)
}

func (s *GearboxServiceImpl) Update(ctx context.Context, id, name, info string) error {
	return s.repo.Update(ctx, &domain.Gearbox{ID: id, GearboxName: name, GearboxInfo: info, UpdatedAt: time.Now()})
}

func (s *GearboxServiceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *GearboxServiceImpl) List(ctx context.Context) ([]domain.Gearbox, error) {
	gearboxes, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Gearbox, len(gearboxes))
	for i, g := range gearboxes {
		result[i] = *g
	}
	return result, nil
}

func (s *GearboxServiceImpl) GetByID(ctx context.Context, id string) (*domain.Gearbox, error) {
	return s.repo.GetByID(ctx, id)
}
