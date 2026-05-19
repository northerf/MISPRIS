package repository

import (
	"context"

	"MISPRIS/internal/domain"

	"github.com/jmoiron/sqlx"
)

type HoClassRolePostgres struct {
	db *sqlx.DB
}

func NewHoClassRolePostgres(db *sqlx.DB) *HoClassRolePostgres {
	return &HoClassRolePostgres{db: db}
}

// List возвращает все записи ho_class_role; если hoClassID не пустой — фильтрует по нему.
func (r *HoClassRolePostgres) List(ctx context.Context, hoClassID string) ([]*domain.HoClassRole, error) {
	var items []*domain.HoClassRole
	var err error
	if hoClassID == "" {
		err = r.db.SelectContext(ctx, &items,
			`SELECT id, ho_class_id, ho_role_id, is_required FROM ho_class_role ORDER BY id`)
	} else {
		err = r.db.SelectContext(ctx, &items,
			`SELECT id, ho_class_id, ho_role_id, is_required FROM ho_class_role
			 WHERE ho_class_id = $1 ORDER BY id`, hoClassID)
	}
	return items, err
}

// ListByClass — псевдоним для List с обязательным фильтром.
func (r *HoClassRolePostgres) ListByClass(ctx context.Context, hoClassID string) ([]*domain.HoClassRole, error) {
	return r.List(ctx, hoClassID)
}

func (r *HoClassRolePostgres) Create(ctx context.Context, hcr *domain.HoClassRole) (string, error) {
	var id string
	query := `INSERT INTO ho_class_role (ho_class_id, ho_role_id, is_required)
	          VALUES ($1, $2, $3)
	          RETURNING id`
	err := r.db.QueryRowContext(ctx, query, hcr.HoClassID, hcr.HoRoleID, hcr.IsRequired).Scan(&id)
	return id, err
}

func (r *HoClassRolePostgres) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM ho_class_role WHERE id = $1`, id)
	return err
}
