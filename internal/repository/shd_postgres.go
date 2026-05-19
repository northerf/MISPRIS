package repository

import (
	"context"
	"time"

	"MISPRIS/internal/domain"

	"github.com/jmoiron/sqlx"
)

type ShdPostgres struct {
	db *sqlx.DB
}

func NewShdPostgres(db *sqlx.DB) *ShdPostgres {
	return &ShdPostgres{db: db}
}

func (r *ShdPostgres) List(ctx context.Context) ([]*domain.SHD, error) {
	var items []*domain.SHD
	query := `SELECT shd_id, name,
	           COALESCE(shd_type,'') AS shd_type,
	           COALESCE(inn,'') AS inn,
	           COALESCE(description,'') AS description,
	           created_at, updated_at
	          FROM shd ORDER BY shd_id`
	err := r.db.SelectContext(ctx, &items, query)
	return items, err
}

func (r *ShdPostgres) GetByID(ctx context.Context, id string) (*domain.SHD, error) {
	var item domain.SHD
	query := `SELECT shd_id, name,
	           COALESCE(shd_type,'') AS shd_type,
	           COALESCE(inn,'') AS inn,
	           COALESCE(description,'') AS description,
	           created_at, updated_at
	          FROM shd WHERE shd_id = $1`
	err := r.db.GetContext(ctx, &item, query, id)
	return &item, err
}

func (r *ShdPostgres) Create(ctx context.Context, s *domain.SHD) (string, error) {
	var id string
	query := `INSERT INTO shd (name, shd_type, inn, description, created_at, updated_at)
	          VALUES ($1, NULLIF($2,''), NULLIF($3,''), NULLIF($4,''), NOW(), NOW())
	          RETURNING shd_id`
	err := r.db.QueryRowContext(ctx, query, s.Name, s.ShdType, s.INN, s.Description).Scan(&id)
	return id, err
}

func (r *ShdPostgres) Update(ctx context.Context, s *domain.SHD) error {
	query := `UPDATE shd SET name=$1, shd_type=NULLIF($2,''), inn=NULLIF($3,''),
	           description=NULLIF($4,''), updated_at=$5
	          WHERE shd_id = $6`
	_, err := r.db.ExecContext(ctx, query, s.Name, s.ShdType, s.INN, s.Description, time.Now(), s.ID)
	return err
}

func (r *ShdPostgres) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM shd WHERE shd_id = $1`, id)
	return err
}
