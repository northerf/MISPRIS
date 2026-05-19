package repository

import (
	"context"
	"time"

	"MISPRIS/internal/domain"

	"github.com/jmoiron/sqlx"
)

type EnumPositionPostgres struct {
	db *sqlx.DB
}

func NewEnumPositionPostgres(db *sqlx.DB) *EnumPositionPostgres {
	return &EnumPositionPostgres{db: db}
}

func (r *EnumPositionPostgres) List(ctx context.Context) ([]*domain.EnumPosition, error) {
	var items []*domain.EnumPosition
	query := `SELECT enum_position_id, enum_class_id, value, order_num, created_at, updated_at
	          FROM enum_position ORDER BY enum_class_id, order_num`
	err := r.db.SelectContext(ctx, &items, query)
	return items, err
}

func (r *EnumPositionPostgres) GetByID(ctx context.Context, id string) (*domain.EnumPosition, error) {
	var item domain.EnumPosition
	query := `SELECT enum_position_id, enum_class_id, value, order_num, created_at, updated_at
	          FROM enum_position WHERE enum_position_id = $1`
	err := r.db.GetContext(ctx, &item, query, id)
	return &item, err
}

func (r *EnumPositionPostgres) Create(ctx context.Context, p *domain.EnumPosition) (string, error) {
	var id string
	query := `INSERT INTO enum_position (enum_class_id, value, order_num, created_at, updated_at)
	          VALUES ($1, $2, $3, NOW(), NOW()) RETURNING enum_position_id`
	err := r.db.QueryRowContext(ctx, query, p.EnumClassID, p.Value, p.OrderNum).Scan(&id)
	return id, err
}

func (r *EnumPositionPostgres) Update(ctx context.Context, p *domain.EnumPosition) error {
	query := `UPDATE enum_position SET value=$1, order_num=$2, updated_at=$3
	          WHERE enum_position_id = $4`
	_, err := r.db.ExecContext(ctx, query, p.Value, p.OrderNum, time.Now(), p.ID)
	return err
}

func (r *EnumPositionPostgres) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM enum_position WHERE enum_position_id = $1`, id)
	return err
}
