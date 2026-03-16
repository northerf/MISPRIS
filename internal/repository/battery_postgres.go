package repository

import (
	"context"

	"MISPRIS/internal/domain"

	"github.com/jmoiron/sqlx"
)

type BatteryPostgres struct {
	db *sqlx.DB
}

func NewBatteryPostgres(db *sqlx.DB) *BatteryPostgres {
	return &BatteryPostgres{db: db}
}

func (r *BatteryPostgres) GetByID(ctx context.Context, id int64) (*domain.Battery, error) {
	var b domain.Battery
	query := `SELECT battery_id, battery_name, battery_type, battery_capacity, battery_info
	          FROM battery WHERE battery_id = $1`
	err := r.db.GetContext(ctx, &b, query, id)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BatteryPostgres) Create(ctx context.Context, b *domain.Battery) (int64, error) {
	var id int64
	query := `INSERT INTO battery (battery_name, battery_type, battery_capacity, battery_info)
	          VALUES ($1, $2, $3, $4) RETURNING battery_id`
	err := r.db.QueryRowContext(ctx, query,
		b.BatteryName, b.BatteryType, b.BatteryCapacity, b.BatteryInfo,
	).Scan(&id)
	return id, err
}

func (r *BatteryPostgres) Update(ctx context.Context, b *domain.Battery) error {
	query := `UPDATE battery
	          SET battery_name = $1, battery_type = $2, battery_capacity = $3, battery_info = $4
	          WHERE battery_id = $5`
	_, err := r.db.ExecContext(ctx, query,
		b.BatteryName, b.BatteryType, b.BatteryCapacity, b.BatteryInfo, b.ID,
	)
	return err
}

func (r *BatteryPostgres) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM battery WHERE battery_id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *BatteryPostgres) List(ctx context.Context) ([]*domain.Battery, error) {
	var batteries []domain.Battery
	query := `SELECT battery_id, battery_name, battery_type, battery_capacity, battery_info FROM battery`
	err := r.db.SelectContext(ctx, &batteries, query)
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Battery, len(batteries))
	for i := range batteries {
		result[i] = &batteries[i]
	}
	return result, nil
}
