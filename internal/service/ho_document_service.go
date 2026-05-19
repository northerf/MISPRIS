package service

import (
	"context"

	"MISPRIS/internal/domain"
	"MISPRIS/internal/repository"
)

type HoDocumentServiceImpl struct {
	repo repository.HoDocumentRepository
}

func NewHoDocumentService(repo repository.HoDocumentRepository) *HoDocumentServiceImpl {
	return &HoDocumentServiceImpl{repo: repo}
}

func (s *HoDocumentServiceImpl) ListByHo(ctx context.Context, hoID string) ([]*domain.HoDocument, error) {
	return s.repo.ListByHo(ctx, hoID)
}

func (s *HoDocumentServiceImpl) Create(ctx context.Context, hoID, docClassID, docNumber, docDate, note string) (string, error) {
	item := &domain.HoDocument{
		HoID:       hoID,
		DocClassID: docClassID,
		DocNumber:  docNumber,
		DocDate:    docDate,
		Note:       note,
	}
	return s.repo.Create(ctx, item)
}

func (s *HoDocumentServiceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
