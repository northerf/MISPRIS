package repository

import (
	"context"
	"fmt"

	"MISPRIS/internal/domain"

	"github.com/jmoiron/sqlx"
)

type ComponentParameterPostgres struct {
	db *sqlx.DB
}

func NewComponentParameterPostgres(db *sqlx.DB) *ComponentParameterPostgres {
	return &ComponentParameterPostgres{db: db}
}

func (r *ComponentParameterPostgres) List(ctx context.Context) ([]*domain.ComponentParameter, error) {
	var items []*domain.ComponentParameter
	query := `SELECT component_parameter_id, component_type, parameter_id,
	           order_num, COALESCE(min_val,0) AS min_val, COALESCE(max_val,0) AS max_val
	          FROM component_parameter ORDER BY component_type, order_num`
	err := r.db.SelectContext(ctx, &items, query)
	return items, err
}

func (r *ComponentParameterPostgres) GetByID(ctx context.Context, id string) (*domain.ComponentParameter, error) {
	var item domain.ComponentParameter
	query := `SELECT component_parameter_id, component_type, parameter_id,
	           order_num, COALESCE(min_val,0) AS min_val, COALESCE(max_val,0) AS max_val
	          FROM component_parameter WHERE component_parameter_id = $1`
	err := r.db.GetContext(ctx, &item, query, id)
	return &item, err
}

func (r *ComponentParameterPostgres) Create(ctx context.Context, cp *domain.ComponentParameter) (string, error) {
	var id string
	query := `INSERT INTO component_parameter (component_type, parameter_id, order_num, min_val, max_val)
	          VALUES ($1, $2, $3, NULLIF($4,0), NULLIF($5,0)) RETURNING component_parameter_id`
	err := r.db.QueryRowContext(ctx, query,
		cp.ComponentType, cp.ParameterID, cp.OrderNum, cp.MinVal, cp.MaxVal,
	).Scan(&id)
	return id, err
}

func (r *ComponentParameterPostgres) Update(ctx context.Context, cp *domain.ComponentParameter) error {
	query := `UPDATE component_parameter SET order_num=$1, min_val=NULLIF($2,0), max_val=NULLIF($3,0)
	          WHERE component_parameter_id = $4`
	_, err := r.db.ExecContext(ctx, query, cp.OrderNum, cp.MinVal, cp.MaxVal, cp.ID)
	return err
}

func (r *ComponentParameterPostgres) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM component_parameter WHERE component_parameter_id = $1`, id)
	return err
}

// GetByType — прямой JOIN с COALESCE вместо SQL-функции (обходим NULL-поля).
func (r *ComponentParameterPostgres) GetByType(ctx context.Context, componentType string) ([]*domain.ComponentParameterFull, error) {
	var items []*domain.ComponentParameterFull
	query := `
		SELECT
			cp.component_parameter_id        AS cp_id,
			p.parameter_id                   AS param_id,
			p.designation,
			p.name,
			p.param_type,
			COALESCE(p.measuring_unit, '') AS measuring_unit,
			COALESCE(cp.min_val, 0)         AS min_val,
			COALESCE(cp.max_val, 0)         AS max_val,
			cp.order_num
		FROM component_parameter cp
		JOIN parameter p ON cp.parameter_id = p.parameter_id
		WHERE cp.component_type = $1
		ORDER BY cp.order_num`
	err := r.db.SelectContext(ctx, &items, query, componentType)
	if err != nil {
		return nil, fmt.Errorf("GetByType: %w", err)
	}
	return items, nil
}

// CopyFromType вызывает SQL-процедуру copy_component_parameters — наследование параметров.
func (r *ComponentParameterPostgres) CopyFromType(ctx context.Context, fromType, toType string) error {
	_, err := r.db.ExecContext(ctx, `CALL copy_component_parameters($1, $2)`, fromType, toType)
	return err
}
