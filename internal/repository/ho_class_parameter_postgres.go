package repository

import (
	"context"
	"fmt"

	"MISPRIS/internal/domain"

	"github.com/jmoiron/sqlx"
)

type HoClassParameterPostgres struct {
	db *sqlx.DB
}

func NewHoClassParameterPostgres(db *sqlx.DB) *HoClassParameterPostgres {
	return &HoClassParameterPostgres{db: db}
}

func (r *HoClassParameterPostgres) List(ctx context.Context, hoClassID string) ([]*domain.HoClassParameter, error) {
	var items []*domain.HoClassParameter
	if hoClassID != "" {
		err := r.db.SelectContext(ctx, &items,
			`SELECT id, ho_class_id, parameter_id, order_num,
			        COALESCE(min_val,0) AS min_val, COALESCE(max_val,0) AS max_val
			 FROM ho_class_parameter WHERE ho_class_id = $1 ORDER BY order_num`, hoClassID)
		return items, err
	}
	err := r.db.SelectContext(ctx, &items,
		`SELECT id, ho_class_id, parameter_id, order_num,
		        COALESCE(min_val,0) AS min_val, COALESCE(max_val,0) AS max_val
		 FROM ho_class_parameter ORDER BY ho_class_id, order_num`)
	return items, err
}

func (r *HoClassParameterPostgres) GetByID(ctx context.Context, id string) (*domain.HoClassParameter, error) {
	var item domain.HoClassParameter
	err := r.db.GetContext(ctx, &item,
		`SELECT id, ho_class_id, parameter_id, order_num,
		        COALESCE(min_val,0) AS min_val, COALESCE(max_val,0) AS max_val
		 FROM ho_class_parameter WHERE id = $1`, id)
	return &item, err
}

func (r *HoClassParameterPostgres) Create(ctx context.Context, cp *domain.HoClassParameter) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO ho_class_parameter (ho_class_id, parameter_id, order_num, min_val, max_val)
		 VALUES ($1,$2,$3,NULLIF($4,0),NULLIF($5,0)) RETURNING id`,
		cp.HoClassID, cp.ParameterID, cp.OrderNum, cp.MinVal, cp.MaxVal,
	).Scan(&id)
	return id, err
}

func (r *HoClassParameterPostgres) Update(ctx context.Context, cp *domain.HoClassParameter) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE ho_class_parameter SET order_num=$1, min_val=NULLIF($2,0), max_val=NULLIF($3,0)
		 WHERE id = $4`,
		cp.OrderNum, cp.MinVal, cp.MaxVal, cp.ID)
	return err
}

func (r *HoClassParameterPostgres) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM ho_class_parameter WHERE id = $1`, id)
	return err
}

// GetByHoClass — прямой JOIN без SQL-функции, COALESCE для NULL-полей
// Заменяет get_ho_class_parameters() чтобы избежать Scan NULL→float64
func (r *HoClassParameterPostgres) GetByHoClass(ctx context.Context, hoClassID string) ([]*domain.HoClassParameterFull, error) {
	var items []*domain.HoClassParameterFull
	query := `
		SELECT
			cp.id                                AS cp_id,
			p.parameter_id                       AS param_id,
			p.designation,
			p.name,
			p.param_type,
			COALESCE(p.measuring_unit,'')        AS measuring_unit,
			COALESCE(cp.min_val, 0)              AS min_val,
			COALESCE(cp.max_val, 0)              AS max_val,
			cp.order_num
		FROM ho_class_parameter cp
		JOIN parameter p ON cp.parameter_id = p.parameter_id
		WHERE cp.ho_class_id = $1
		ORDER BY cp.order_num`
	err := r.db.SelectContext(ctx, &items, query, hoClassID)
	if err != nil {
		return nil, fmt.Errorf("GetByHoClass: %w", err)
	}
	return items, nil
}

func (r *HoClassParameterPostgres) CopyFromClass(ctx context.Context, fromClassID, toClassID string) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO ho_class_parameter (ho_class_id, parameter_id, order_num, min_val, max_val)
		 SELECT $2, parameter_id, order_num, min_val, max_val
		 FROM ho_class_parameter WHERE ho_class_id = $1
		 ON CONFLICT (ho_class_id, parameter_id) DO NOTHING`,
		fromClassID, toClassID)
	return err
}
