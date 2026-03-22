package service

import (
	"context"
	"fmt"

	"MISPRIS/internal/domain"
	"MISPRIS/internal/repository"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

//Charger
type ChargerServiceImpl struct {
	repo repository.ChargerRepository
}

func NewChargerService(repo repository.ChargerRepository) *ChargerServiceImpl {
	return &ChargerServiceImpl{repo: repo}
}

func (s *ChargerServiceImpl) Create(ctx context.Context, name, info string) (string, error) {
	return s.repo.Create(ctx, &domain.Charger{
		ID:   uuid.New().String(),
		Name: name,
		Info: info,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: false,
	})
}

func (s *ChargerServiceImpl) Update(ctx context.Context, id, name, info string) error {
	return s.repo.Update(ctx, &domain.Charger{ID: id, Name: name, Info: info})
}

func (s *ChargerServiceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *ChargerServiceImpl) List(ctx context.Context) ([]domain.Charger, error) {
	rows, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Charger, len(rows))
	for i, r := range rows {
		result[i] = *r
	}
	return result, nil
}

func (s *ChargerServiceImpl) GetByID(ctx context.Context, id string) (*domain.Charger, error) {
	return s.repo.GetByID(ctx, id)
}

//Connector
type ConnectorServiceImpl struct {
	repo repository.ConnectorRepository
}

func NewConnectorService(repo repository.ConnectorRepository) *ConnectorServiceImpl {
	return &ConnectorServiceImpl{repo: repo}
}

func (s *ConnectorServiceImpl) Create(ctx context.Context, name, info string) (string, error) {
	return s.repo.Create(ctx, &domain.Connector{
		ID:   uuid.New().String(),
		Name: name,
		Info: info,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		DeletedAt: false,
	})
}

func (s *ConnectorServiceImpl) Update(ctx context.Context, id, name, info string) error {
	return s.repo.Update(ctx, &domain.Connector{ID: id, Name: name, Info: info})
}

func (s *ConnectorServiceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *ConnectorServiceImpl) List(ctx context.Context) ([]domain.Connector, error) {
	rows, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Connector, len(rows))
	for i, r := range rows {
		result[i] = *r
	}
	return result, nil
}

func (s *ConnectorServiceImpl) GetByID(ctx context.Context, id string) (*domain.Connector, error) {
	return s.repo.GetByID(ctx, id)
}

//ChargerSystem
type ChargerSystemServiceImpl struct {
	db               *sqlx.DB
	chargerSystemRepo repository.ChargerSystemRepository
	chargerService    ChargerService
	connectorService  ConnectorService
}

func NewChargerSystemService(
	db *sqlx.DB,
	chargerSystemRepo repository.ChargerSystemRepository,
	chargerService ChargerService,
	connectorService ConnectorService,
) *ChargerSystemServiceImpl {
	return &ChargerSystemServiceImpl{
		db:                db,
		chargerSystemRepo: chargerSystemRepo,
		chargerService:    chargerService,
		connectorService:  connectorService,
	}
}

func (s *ChargerSystemServiceImpl) Create(ctx context.Context, chargerID, connectorID string) (string, error) {
	
	if _, err := s.chargerService.GetByID(ctx, chargerID); err != nil {
		return "", fmt.Errorf("charger %s not found: %w", chargerID, err)
	}
	
	if _, err := s.connectorService.GetByID(ctx, connectorID); err != nil {
		return "", fmt.Errorf("connector %s not found: %w", connectorID, err)
	}
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	id, err := s.chargerSystemRepo.CreateTx(ctx, tx, &domain.ChargerSystem{
		ChargerID:   chargerID,
		ConnectorID: connectorID,
	})
	if err != nil {
		return "", err
	}
	return id, tx.Commit()
}

func (s *ChargerSystemServiceImpl) Update(ctx context.Context, id, chargerID, connectorID string) error {
	return s.chargerSystemRepo.Update(ctx, &domain.ChargerSystem{
		ID:          id,
		ChargerID:   chargerID,
		ConnectorID: connectorID,
	})
}

func (s *ChargerSystemServiceImpl) Delete(ctx context.Context, id string) error {
	return s.chargerSystemRepo.Delete(ctx, id)
}

func (s *ChargerSystemServiceImpl) List(ctx context.Context) ([]domain.ChargerSystem, error) {
	rows, err := s.chargerSystemRepo.List(ctx, id)
	if err != nil {
		return nil, err
	}
	result := make([]domain.ChargerSystem, len(rows))
	for i, r := range rows {
		result[i] = *r
	}
	return result, nil
}

func (s *ChargerSystemServiceImpl) GetByID(ctx context.Context, id string) (*domain.ChargerSystem, error) {
	return s.chargerSystemRepo.GetByID(ctx, id)
}