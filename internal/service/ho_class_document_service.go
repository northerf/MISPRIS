package service

import (
	"context"

	"MISPRIS/internal/domain"
	"MISPRIS/internal/repository"
)

type HoClassDocumentServiceImpl struct {
	repo repository.HoClassDocumentRepository
}

func NewHoClassDocumentService(repo repository.HoClassDocumentRepository) *HoClassDocumentServiceImpl {
	return &HoClassDocumentServiceImpl{repo: repo}
}

func (s *HoClassDocumentServiceImpl) ListByClass(ctx context.Context, hoClassID string) ([]*domain.HoClassDocument, error) {
	return s.repo.ListByClass(ctx, hoClassID)
}

func (s *HoClassDocumentServiceImpl) Create(ctx context.Context, hoClassID, docClassID, roleName string, isRequired bool) (string, error) {
	item := &domain.HoClassDocument{
		HoClassID:  hoClassID,
		DocClassID: docClassID,
		RoleName:   roleName,
		IsRequired: isRequired,
	}
	return s.repo.Create(ctx, item)
}

func (s *HoClassDocumentServiceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
