package service

import (
	"context"
	"fmt"

	"MISPRIS/internal/domain"
	"MISPRIS/internal/repository"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

//Frame
type FrameServiceImpl struct {
	repo repository.FrameRepository
}

func NewFrameService(repo repository.FrameRepository) *FrameServiceImpl {
	return &FrameServiceImpl{repo: repo}
}

func (s *FrameServiceImpl) Create(ctx context.Context, name, info string) (string, error) {
	return s.repo.Create(ctx, &domain.Frame{
		ID:   uuid.New().String(),
		Name: name,
		Info: info,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: false,
	})
}

func (s *FrameServiceImpl) Update(ctx context.Context, id, name, info string) error {
	return s.repo.Update(ctx, &domain.Frame{ID: id, Name: name, Info: info})
}

func (s *FrameServiceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *FrameServiceImpl) List(ctx context.Context) ([]domain.Frame, error) {
	rows, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Frame, len(rows))
	for i, r := range rows {
		result[i] = *r
	}
	return result, nil
}

func (s *FrameServiceImpl) GetByID(ctx context.Context, id string) (*domain.Frame, error) {
	return s.repo.GetByID(ctx, id)
}

//Suspension
type SuspensionServiceImpl struct {
	repo repository.SuspensionRepository
}

func NewSuspensionService(repo repository.SuspensionRepository) *SuspensionServiceImpl {
	return &SuspensionServiceImpl{repo: repo}
}

func (s *SuspensionServiceImpl) Create(ctx context.Context, name, info string) (string, error) {
	return s.repo.Create(ctx, &domain.Suspension{
		ID:   uuid.New().String(),
		Name: name,
		Info: info,
	})
}

func (s *SuspensionServiceImpl) Update(ctx context.Context, id, name, info string) error {
	return s.repo.Update(ctx, &domain.Suspension{ID: id, Name: name, Info: info})
}

func (s *SuspensionServiceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *SuspensionServiceImpl) List(ctx context.Context) ([]domain.Suspension, error) {
	rows, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Suspension, len(rows))
	for i, r := range rows {
		result[i] = *r
	}
	return result, nil
}

func (s *SuspensionServiceImpl) GetByID(ctx context.Context, id string) (*domain.Suspension, error) {
	return s.repo.GetByID(ctx, id)
}

//BreakSystem
type BreakSystemServiceImpl struct {
	repo repository.BreakSystemRepository
}

func NewBreakSystemService(repo repository.BreakSystemRepository) *BreakSystemServiceImpl {
	return &BreakSystemServiceImpl{repo: repo}
}

func (s *BreakSystemServiceImpl) Create(ctx context.Context, name, info string) (string, error) {
	return s.repo.Create(ctx, &domain.BreakSystem{
		ID:   uuid.New().String(),
		Name: name,
		Info: info,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: false,
	})
}

func (s *BreakSystemServiceImpl) Update(ctx context.Context, id, name, info string) error {
	return s.repo.Update(ctx, &domain.BreakSystem{ID: id, Name: name, Info: info})
}

func (s *BreakSystemServiceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *BreakSystemServiceImpl) List(ctx context.Context) ([]domain.BreakSystem, error) {
	rows, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.BreakSystem, len(rows))
	for i, r := range rows {
		result[i] = *r
	}
	return result, nil
}

func (s *BreakSystemServiceImpl) GetByID(ctx context.Context, id string) (*domain.BreakSystem, error) {
	return s.repo.GetByID(ctx, id)
}

//Chassis 
type ChassisServiceImpl struct {
	db             *sqlx.DB
	chassisRepo    repository.ChassisRepository
	frameService   FrameService
	suspensionService SuspensionService
	breakSystemService BreakSystemService
}

func NewChassisService(
	db *sqlx.DB,
	chassisRepo repository.ChassisRepository,
	frameService FrameService,
	suspensionService SuspensionService,
	breakSystemService BreakSystemService,
) *ChassisServiceImpl {
	return &ChassisServiceImpl{
		db:                 db,
		chassisRepo:        chassisRepo,
		frameService:       frameService,
		suspensionService:  suspensionService,
		breakSystemService: breakSystemService,
	}
}

func (s *ChassisServiceImpl) Create(ctx context.Context, frameID, suspensionID, breakSystemID string) (string, error) {
	if _, err := s.frameService.GetByID(ctx, frameID); err != nil {
		return "", fmt.Errorf("frame %s not found: %w", frameID, err)
	}

	if _, err := s.suspensionService.GetByID(ctx, suspensionID); err != nil {
		return "", fmt.Errorf("suspension %s not found: %w", suspensionID, err)
	}
	
	if _, err := s.breakSystemService.GetByID(ctx, breakSystemID); err != nil {
		return "", fmt.Errorf("break system %s not found: %w", breakSystemID, err)
	}
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	id, err := s.chassisRepo.CreateTx(ctx, tx, &domain.Chassis{
		FrameID:       frameID,
		SuspensionID:  suspensionID,
		BreakSystemID: breakSystemID,
	})
	if err != nil {
		return "", err
	}
	return id, tx.Commit()
}

func (s *ChassisServiceImpl) Update(ctx context.Context, id, frameID, suspensionID, breakSystemID string) error {
	return s.chassisRepo.Update(ctx, &domain.Chassis{
		ID:            id,
		FrameID:       frameID,
		SuspensionID:  suspensionID,
		BreakSystemID: breakSystemID,
	})
}

func (s *ChassisServiceImpl) Delete(ctx context.Context, id string) error {
	return s.chassisRepo.Delete(ctx, id)
}

func (s *ChassisServiceImpl) List(ctx context.Context) ([]domain.Chassis, error) {
	rows, err := s.chassisRepo.List(ctx, id)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Chassis, len(rows))
	for i, r := range rows {
		result[i] = *r
	}
	return result, nil
}

func (s *ChassisServiceImpl) GetByID(ctx context.Context, id string) (*domain.Chassis, error) {
	return s.chassisRepo.GetByID(ctx, id)
}