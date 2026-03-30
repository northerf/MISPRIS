package service

import (
	"context"
	"errors"
	"time"

	"MISPRIS/internal/domain"
	"MISPRIS/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{repo: repo}
}

func (s *UserServiceImpl) Create(ctx context.Context, username, password, role string, isActive bool) (string, error) {
	existing, _ := s.repo.GetByUsername(ctx, username)
	if existing != nil {
		return "", errors.New("username already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := &domain.User{
		ID:        uuid.New().String(),
		Username:  username,
		Password:  string(hashed),
		Role:      role,
		IsActive:  isActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return s.repo.Create(ctx, user)
}

func (s *UserServiceImpl) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserServiceImpl) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	return s.repo.GetByUsername(ctx, username)
}

func (s *UserServiceImpl) Update(ctx context.Context, id, username, password, role string, isActive bool) error {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	user.Username = username
	user.Role = role
	user.IsActive = isActive
	user.UpdatedAt = time.Now()

	if password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashed)
	}

	return s.repo.Update(ctx, user)
}

func (s *UserServiceImpl) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *UserServiceImpl) List(ctx context.Context) ([]domain.User, error) {
	users, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.User, len(users))
	for i, u := range users {
		result[i] = *u
	}
	return result, nil
}

func (s *UserServiceImpl) Authenticate(ctx context.Context, username, password string) (*domain.User, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil || user == nil {
		return nil, errors.New("invalid credentials")
	}
	if !user.IsActive {
		return nil, errors.New("user is inactive")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}
