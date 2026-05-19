package repository

import (
	"context"

	"MISPRIS/internal/domain"

	"github.com/jmoiron/sqlx"
)

type HoRolePostgres struct {
	db *sqlx.DB
}

func NewHoRolePostgres(db *sqlx.DB) *HoRolePostgres {
	return &HoRolePostgres{db: db}
}

func (r *HoRolePostgres) List(ctx context.Context) ([]*domain.HoRole, error) {
	var items []*domain.HoRole
	query := `SELECT ho_role_id, name,
	           COALESCE(description,'') AS description
	          FROM ho_role ORDER BY ho_role_id`
	err := r.db.SelectContext(ctx, &items, query)
	return items, err
}

func (r *HoRolePostgres) GetByID(ctx context.Context, id string) (*domain.HoRole, error) {
	var item domain.HoRole
	query := `SELECT ho_role_id, name,
	           COALESCE(description,'') AS description
	          FROM ho_role WHERE ho_role_id = $1`
	err := r.db.GetContext(ctx, &item, query, id)
	return &item, err
}

func (r *HoRolePostgres) Create(ctx context.Context, hr *domain.HoRole) (string, error) {
	var id string
	query := `INSERT INTO ho_role (name, description)
	          VALUES ($1, NULLIF($2,''))
	          RETURNING ho_role_id`
	err := r.db.QueryRowContext(ctx, query, hr.Name, hr.Description).Scan(&id)
	return id, err
}

func (r *HoRolePostgres) Update(ctx context.Context, hr *domain.HoRole) error {
	query := `UPDATE ho_role SET name=$1, description=NULLIF($2,'')
	          WHERE ho_role_id = $3`
	_, err := r.db.ExecContext(ctx, query, hr.Name, hr.Description, hr.ID)
	return err
}

func (r *HoRolePostgres) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM ho_role WHERE ho_role_id = $1`, id)
	return err
}
