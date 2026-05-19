package service

import (
	"context"
	"fmt"
	"strconv"

	"MISPRIS/internal/domain"
	"MISPRIS/internal/repository"
)

type EnumPositionServiceImpl struct {
	repo repository.EnumPositionRepository
}

func NewEnumPositionService(repo repository.EnumPositionRepository) *EnumPositionServiceImpl {
	return &EnumPositionServiceImpl{repo: repo}
}

func (s *EnumPositionServiceImpl) List(ctx context.Context) ([]*domain.EnumPosition, error) {
	return s.repo.List(ctx)
}

func (s *EnumPositionServiceImpl) GetByID(ctx context.Context, id string) (*domain.EnumPosition, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *EnumPositionServiceImpl) Create(ctx context.Context, enumClassID, value, orderNumStr string) (string, error) {
	if enumClassID == "" || value == "" {
		return "", fmt.Errorf("enum_class_id and value are required")
	}
	orderNum := 0
	if orderNumStr != "" {
		var err error
		orderNum, err = strconv.Atoi(orderNumStr)
		if err != nil {
			return "", fmt.Errorf("invalid order_num: %s", orderNumStr)
		}
	}
	p := &domain.EnumPosition{
		EnumClassID: enumClassID,
		Value:       value,
		OrderNum:    orderNum,
	}
	return s.repo.Create(ctx, p)
}

func (s *EnumPositionServiceImpl) Update(ctx context.Context, id, value, orderNumStr string) error {
	orderNum := 0
	if orderNumStr != "" {
		var err error
		orderNum, err = strconv.Atoi(orderNumStr)
		if err != nil {
			return fmt.Errorf("invalid order_num: %s", orderNumStr)
		}
	}
	p := &domain.EnumPosition{
		ID:       id,
		Value:    value,
		OrderNum: orderNum,
	}
	return s.repo.Update(ctx, p)
}

func (s *EnumPositionServiceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
