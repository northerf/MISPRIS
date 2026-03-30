package repository

import (
	"context"
	"database/sql"

	"MISPRIS/internal/domain"

	"github.com/jmoiron/sqlx"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) Create(ctx context.Context, user *domain.User) (string, error) {
	var id string
	query := `INSERT INTO users (user_id, username, password, role, is_active, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING user_id`
	err := r.db.QueryRowContext(ctx, query,
		user.ID, user.Username, user.Password, user.Role,
		user.IsActive, user.CreatedAt, user.UpdatedAt,
	).Scan(&id)
	return id, err
}

func (r *UserPostgres) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	query := `SELECT user_id, username, password, role, is_active, created_at, updated_at
	          FROM users WHERE user_id = $1`
	err := r.db.GetContext(ctx, &user, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserPostgres) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	query := `SELECT user_id, username, password, role, is_active, created_at, updated_at
	          FROM users WHERE username = $1`
	err := r.db.GetContext(ctx, &user, query, username)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserPostgres) Update(ctx context.Context, user *domain.User) error {
	query := `UPDATE users SET username=$1, password=$2, role=$3, is_active=$4, updated_at=$5
	          WHERE user_id=$6`
	_, err := r.db.ExecContext(ctx, query,
		user.Username, user.Password, user.Role, user.IsActive, user.UpdatedAt, user.ID)
	return err
}

func (r *UserPostgres) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE user_id=$1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *UserPostgres) List(ctx context.Context) ([]*domain.User, error) {
	var rows []domain.User
	query := `SELECT user_id, username, password, role, is_active, created_at, updated_at
	          FROM users`
	err := r.db.SelectContext(ctx, &rows, query)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.User, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}
	return result, nil
}
