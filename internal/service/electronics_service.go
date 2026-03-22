package service

import (
	"MISPRIS/internal/domain"
	"MISPRIS/internal/repository"
	"context"
	"fmt"

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
		Name: controllerName,
		Info: controllerInfo,
	})
}
