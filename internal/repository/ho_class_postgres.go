package repository

import (
	"context"
	"time"

	"MISPRIS/internal/domain"

	"github.com/jmoiron/sqlx"
)

type HoClassPostgres struct {
	db *sqlx.DB
}

func NewHoClassPostgres(db *sqlx.DB) *HoClassPostgres {
	return &HoClassPostgres{db: db}
}

func (r *HoClassPostgres) List(ctx context.Context) ([]*domain.HoClass, error) {
	var items []*domain.HoClass
	query := `SELECT ho_class_id, name,
	           COALESCE(designation,'') AS designation,
	           COALESCE(parent_id::text,'') AS parent_id,
	           is_terminal, created_at, updated_at
	          FROM ho_class ORDER BY ho_class_id`
	err := r.db.SelectContext(ctx, &items, query)
	return items, err
}

func (r *HoClassPostgres) GetByID(ctx context.Context, id string) (*domain.HoClass, error) {
	var item domain.HoClass
	query := `SELECT ho_class_id, name,
	           COALESCE(designation,'') AS designation,
	           COALESCE(parent_id::text,'') AS parent_id,
	           is_terminal, created_at, updated_at
	          FROM ho_class WHERE ho_class_id = $1`
	err := r.db.GetContext(ctx, &item, query, id)
	return &item, err
}

func (r *HoClassPostgres) Create(ctx context.Context, h *domain.HoClass) (string, error) {
	var id string
	query := `INSERT INTO ho_class (name, designation, parent_id, is_terminal, created_at, updated_at)
	          VALUES ($1, NULLIF($2,''), NULLIF($3,'')::int, $4, NOW(), NOW())
	          RETURNING ho_class_id`
	err := r.db.QueryRowContext(ctx, query, h.Name, h.Designation, h.ParentID, h.IsTerminal).Scan(&id)
	return id, err
}

func (r *HoClassPostgres) Update(ctx context.Context, h *domain.HoClass) error {
	query := `UPDATE ho_class SET name=$1, designation=NULLIF($2,''),
	           parent_id=NULLIF($3,'')::int, is_terminal=$4, updated_at=$5
	          WHERE ho_class_id = $6`
	_, err := r.db.ExecContext(ctx, query, h.Name, h.Designation, h.ParentID, h.IsTerminal, time.Now(), h.ID)
	return err
}

func (r *HoClassPostgres) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM ho_class WHERE ho_class_id = $1`, id)
	return err
}

// GetTerminal возвращает конечные (терминальные) узлы классификатора ХО.
func (r *HoClassPostgres) GetTerminal(ctx context.Context) ([]*domain.HoClass, error) {
	var items []*domain.HoClass
	query := `SELECT ho_class_id, name,
	           COALESCE(designation,'') AS designation,
	           COALESCE(parent_id::text,'') AS parent_id,
	           is_terminal, created_at, updated_at
	          FROM ho_class WHERE is_terminal = TRUE ORDER BY ho_class_id`
	err := r.db.SelectContext(ctx, &items, query)
	return items, err
}

// GetChildren возвращает дочерние узлы по parentID.
func (r *HoClassPostgres) GetChildren(ctx context.Context, parentID string) ([]*domain.HoClass, error) {
	var items []*domain.HoClass
	query := `SELECT ho_class_id, name,
	           COALESCE(designation,'') AS designation,
	           COALESCE(parent_id::text,'') AS parent_id,
	           is_terminal, created_at, updated_at
	          FROM ho_class WHERE parent_id = $1 ORDER BY ho_class_id`
	err := r.db.SelectContext(ctx, &items, query, parentID)
	return items, err
}
