package service

import (
	"MISPRIS/internal/domain"
	"MISPRIS/internal/repository"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

//Carcass

type CarcassServiceImpl struct {
	repo repository.CarcassRepository
}

func NewCarcassService(repo repository.CarcassRepository) *CarcassServiceImpl {
	return &CarcassServiceImpl{repo: repo}
}

func (s *CarcassServiceImpl) Create(ctx context.Context, name, info string) (string, error) {
	return s.repo.Create(ctx, &domain.Carcass{
		ID:   uuid.New().String(),
		Name: name,
		Info: info,
	})
}

func (s *CarcassServiceImpl) Update(ctx context.Context, id, name, info string) error {
	return s.repo.Update(ctx, &domain.Carcass{ID: id, Name: name, Info: info})
}

func (s *CarcassServiceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *CarcassServiceImpl) List(ctx context.Context) ([]domain.Carcass, error) {
	rows, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Carcass, len(rows))
	for i, r := range rows {
		result[i] = *r
	}
	return result, nil
}

func (s *CarcassServiceImpl) GetByID(ctx context.Context, id string) (*domain.Carcass, error) {
	return s.repo.GetByID(ctx, id)
}

//Doors

type DoorsServiceImpl struct {
	repo repository.DoorsRepository
}

func NewDoorsService(repo repository.DoorsRepository) *DoorsServiceImpl {
	return &DoorsServiceImpl{repo: repo}
}

func (s *DoorsServiceImpl) Create(ctx context.Context, name, info string) (string, error) {
	return s.repo.Create(ctx, &domain.Doors{
		ID:   uuid.New().String(),
		Name: name,
		Info: info,
	})
}

func (s *DoorsServiceImpl) Update(ctx context.Context, id, name, info string) error {
	return s.repo.Update(ctx, &domain.Doors{ID: id, Name: name, Info: info})
}

func (s *DoorsServiceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *DoorsServiceImpl) List(ctx context.Context) ([]domain.Doors, error) {
	rows, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Doors, len(rows))
	for i, r := range rows {
		result[i] = *r
	}
	return result, nil
}

func (s *DoorsServiceImpl) GetByID(ctx context.Context, id string) (*domain.Doors, error) {
	return s.repo.GetByID(ctx, id)
}

//Wings

type WingsServiceImpl struct {
	repo repository.WingsRepository
}

func NewWingsService(repo repository.WingsRepository) *WingsServiceImpl {
	return &WingsServiceImpl{repo: repo}
}

func (s *WingsServiceImpl) Create(ctx context.Context, name, info string) (string, error) {
	return s.repo.Create(ctx, &domain.Wings{
		ID:   uuid.New().String(),
		Name: name,
		Info: info,
	})
}

func (s *WingsServiceImpl) Update(ctx context.Context, id, name, info string) error {
	return s.repo.Update(ctx, &domain.Wings{ID: id, Name: name, Info: info})
}

func (s *WingsServiceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *WingsServiceImpl) List(ctx context.Context) ([]domain.Wings, error) {
	rows, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Wings, len(rows))
	for i, r := range rows {
		result[i] = *r
	}
	return result, nil
}

func (s *WingsServiceImpl) GetByID(ctx context.Context, id string) (*domain.Wings, error) {
	return s.repo.GetByID(ctx, id)
}

//Body

type BodyServiceImpl struct {
	db             *sqlx.DB
	bodyRepo       repository.BodyRepository
	carcassService CarcassService
	doorsService   DoorsService
	wingsService   WingsService
}

func NewBodyService(
	db *sqlx.DB,
	bodyRepo repository.BodyRepository,
	carcassService CarcassService,
	doorsService DoorsService,
	wingsService WingsService,
) *BodyServiceImpl {
	return &BodyServiceImpl{
		db:             db,
		bodyRepo:       bodyRepo,
		carcassService: carcassService,
		doorsService:   doorsService,
		wingsService:   wingsService,
	}
}

func (s *BodyServiceImpl) Create(ctx context.Context, carcassID, doorsID, wingsID string) (string, error) {
	if _, err := s.carcassService.GetByID(ctx, carcassID); err != nil {
		return "", fmt.Errorf("carcass %s not found: %w", carcassID, err)
	}
	if _, err := s.doorsService.GetByID(ctx, doorsID); err != nil {
		return "", fmt.Errorf("doors %s not found: %w", doorsID, err)
	}
	if _, err := s.wingsService.GetByID(ctx, wingsID); err != nil {
		return "", fmt.Errorf("wings %s not found: %w", wingsID, err)
	}
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	id, err := s.bodyRepo.CreateTx(ctx, tx, &domain.Body{
		CarcassID: carcassID,
		DoorsID:   doorsID,
		WingsID:   wingsID,
	})
	if err != nil {
		return "", err
	}
	return id, tx.Commit()
}

func (s *BodyServiceImpl) Update(ctx context.Context, id, carcassID, doorsID, wingsID string) error {
	return s.bodyRepo.Update(ctx, &domain.Body{
		ID:        id,
		CarcassID: carcassID,
		DoorsID:   doorsID,
		WingsID:   wingsID,
	})
}

func (s *BodyServiceImpl) Delete(ctx context.Context, id string) error {
	return s.bodyRepo.Delete(ctx, id)
}

func (s *BodyServiceImpl) List(ctx context.Context) ([]domain.Body, error) {
	rows, err := s.bodyRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Body, len(rows))
	for i, r := range rows {
		result[i] = *r
	}
	return result, nil
}

func (s *BodyServiceImpl) GetByID(ctx context.Context, id string) (*domain.Body, error) {
	return s.bodyRepo.GetByID(ctx, id)
}
