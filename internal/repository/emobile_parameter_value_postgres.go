package repository

import (
	"context"

	"MISPRIS/internal/domain"

	"github.com/jmoiron/sqlx"
)

type EmobileParameterValuePostgres struct {
	db *sqlx.DB
}

func NewEmobileParameterValuePostgres(db *sqlx.DB) *EmobileParameterValuePostgres {
	return &EmobileParameterValuePostgres{db: db}
}

func (r *EmobileParameterValuePostgres) List(ctx context.Context) ([]*domain.EmobileParameterValue, error) {
	var items []*domain.EmobileParameterValue
	query := `SELECT value_id, emobile_id, component_parameter_id,
	           COALESCE(val_real,0) AS val_real, COALESCE(val_int,0) AS val_int,
	           COALESCE(val_str,'') AS val_str, COALESCE(enum_val_id::text,'') AS enum_val_id
	          FROM emobile_parameter_value ORDER BY value_id`
	err := r.db.SelectContext(ctx, &items, query)
	return items, err
}

func (r *EmobileParameterValuePostgres) GetByID(ctx context.Context, id string) (*domain.EmobileParameterValue, error) {
	var item domain.EmobileParameterValue
	query := `SELECT value_id, emobile_id, component_parameter_id,
	           COALESCE(val_real,0) AS val_real, COALESCE(val_int,0) AS val_int,
	           COALESCE(val_str,'') AS val_str, COALESCE(enum_val_id::text,'') AS enum_val_id
	          FROM emobile_parameter_value WHERE value_id = $1`
	err := r.db.GetContext(ctx, &item, query, id)
	return &item, err
}

func (r *EmobileParameterValuePostgres) Create(ctx context.Context, v *domain.EmobileParameterValue) (string, error) {
	var id string
	query := `INSERT INTO emobile_parameter_value
	           (emobile_id, component_parameter_id, val_real, val_int, val_str, enum_val_id)
	          VALUES ($1, $2, NULLIF($3,0), NULLIF($4,0), NULLIF($5,''), NULLIF($6,'')::int)
	          RETURNING value_id`
	err := r.db.QueryRowContext(ctx, query,
		v.EmobileID, v.ComponentParameterID,
		v.ValReal, v.ValInt, v.ValStr, v.EnumValID,
	).Scan(&id)
	return id, err
}

func (r *EmobileParameterValuePostgres) Update(ctx context.Context, v *domain.EmobileParameterValue) error {
	query := `UPDATE emobile_parameter_value
	          SET val_real=NULLIF($1,0), val_int=NULLIF($2,0), val_str=NULLIF($3,''),
	              enum_val_id=NULLIF($4,'')::int
	          WHERE value_id = $5`
	_, err := r.db.ExecContext(ctx, query, v.ValReal, v.ValInt, v.ValStr, v.EnumValID, v.ID)
	return err
}

func (r *EmobileParameterValuePostgres) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM emobile_parameter_value WHERE value_id = $1`, id)
	return err
}

// GetByEmobile возвращает все значения параметров конкретного автомобиля.
func (r *EmobileParameterValuePostgres) GetByEmobile(ctx context.Context, emobileID string) ([]*domain.EmobileParameterValue, error) {
	var items []*domain.EmobileParameterValue
	query := `SELECT value_id, emobile_id, component_parameter_id,
	           COALESCE(val_real,0) AS val_real, COALESCE(val_int,0) AS val_int,
	           COALESCE(val_str,'') AS val_str, COALESCE(enum_val_id::text,'') AS enum_val_id
	          FROM emobile_parameter_value WHERE emobile_id = $1`
	err := r.db.SelectContext(ctx, &items, query, emobileID)
	return items, err
}
