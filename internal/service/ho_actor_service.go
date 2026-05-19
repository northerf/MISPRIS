package service

import (
	"context"

	"MISPRIS/internal/domain"
	"MISPRIS/internal/repository"
)

type HoActorServiceImpl struct {
	repo repository.HoActorRepository
}

func NewHoActorService(repo repository.HoActorRepository) *HoActorServiceImpl {
	return &HoActorServiceImpl{repo: repo}
}

func (s *HoActorServiceImpl) ListByHo(ctx context.Context, hoID string) ([]*domain.HoActor, error) {
	return s.repo.ListByHo(ctx, hoID)
}

// Create вызывает SQL-процедуру set_ho_actor для назначения актора с валидацией.
func (s *HoActorServiceImpl) Create(ctx context.Context, hoID, hoRoleID, shdID string) (string, error) {
	item := &domain.HoActor{
		HoID:     hoID,
		HoRoleID: hoRoleID,
		ShdID:    shdID,
	}
	return s.repo.Create(ctx, item)
}

func (s *HoActorServiceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
