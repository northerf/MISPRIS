package repository

import (
	"context"
	"time"

	"MISPRIS/internal/domain"

	"github.com/jmoiron/sqlx"
)

type ParameterPostgres struct {
	db *sqlx.DB
}

func NewParameterPostgres(db *sqlx.DB) *ParameterPostgres {
	return &ParameterPostgres{db: db}
}

func (r *ParameterPostgres) List(ctx context.Context) ([]*domain.Parameter, error) {
	var items []*domain.Parameter
	query := `SELECT parameter_id, designation, name, param_type,
	           COALESCE(measuring_unit,'') AS measuring_unit,
	           COALESCE(enum_class_id::text,'') AS enum_class_id,
	           created_at, updated_at
	          FROM parameter ORDER BY parameter_id`
	err := r.db.SelectContext(ctx, &items, query)
	return items, err
}

func (r *ParameterPostgres) GetByID(ctx context.Context, id string) (*domain.Parameter, error) {
	var item domain.Parameter
	query := `SELECT parameter_id, designation, name, param_type,
	           COALESCE(measuring_unit,'') AS measuring_unit,
	           COALESCE(enum_class_id::text,'') AS enum_class_id,
	           created_at, updated_at
	          FROM parameter WHERE parameter_id = $1`
	err := r.db.GetContext(ctx, &item, query, id)
	return &item, err
}

func (r *ParameterPostgres) Create(ctx context.Context, p *domain.Parameter) (string, error) {
	var id string
	query := `INSERT INTO parameter (designation, name, param_type, measuring_unit, enum_class_id, created_at, updated_at)
	          VALUES ($1, $2, $3, NULLIF($4,''), NULLIF($5,'')::int, NOW(), NOW()) RETURNING parameter_id`
	err := r.db.QueryRowContext(ctx, query, p.Designation, p.Name, p.ParamType, p.MeasuringUnit, p.EnumClassID).Scan(&id)
	return id, err
}

func (r *ParameterPostgres) Update(ctx context.Context, p *domain.Parameter) error {
	query := `UPDATE parameter SET designation=$1, name=$2, param_type=$3,
	           measuring_unit=NULLIF($4,''), enum_class_id=NULLIF($5,'')::int, updated_at=$6
	          WHERE parameter_id = $7`
	_, err := r.db.ExecContext(ctx, query,
		p.Designation, p.Name, p.ParamType, p.MeasuringUnit, p.EnumClassID, time.Now(), p.ID)
	return err
}

func (r *ParameterPostgres) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM parameter WHERE parameter_id = $1`, id)
	return err
}
