package repository

import (
	"context"

	"MISPRIS/internal/domain"

	"github.com/jmoiron/sqlx"
)

type EmobilePostgres struct {
	db *sqlx.DB
}

func NewEmobilePostgres(db *sqlx.DB) *EmobilePostgres {
	return &EmobilePostgres{
		db: db,
	}
}

func (r *EmobilePostgres) GetByID(ctx context.Context, id string) (*domain.Emobile, error) {
	var emobile domain.Emobile

	if err := r.db.GetContext(ctx, &emobile,
		`SELECT * FROM emobile WHERE emobile_id = $1`, id,
	); err != nil {
		return nil, err
	}
	return &emobile, nil
}

func (r *EmobilePostgres) Create(ctx context.Context, tx *sqlx.Tx, emobile *domain.Emobile) (string, error) {
	var emobileID string
	err := tx.QueryRowContext(ctx, `INSERT INTO emobile (emobile_name,
                     power_point_id, battery_id, charger_system_id,
                     chassis_id, body_id, electronics_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING emobile_id`,
		emobile.Name, emobile.PowerPointID, emobile.BatteryID,
		emobile.ChargerSystemID, emobile.ChassisID, emobile.BodyID, emobile.ElectronicsID).Scan(&emobileID)
	return emobileID, err
}

func (r *EmobilePostgres) Update(ctx context.Context, emobile *domain.Emobile) error {
	query := `UPDATE emobile
	          SET emobile_name=$1,
	              power_point_id=$2,
	              battery_id=$3,
	              charger_system_id=$4,
	              chassis_id=$5,
	              body_id=$6,
	              electronics_id=$7
	          WHERE emobile_id=$8`
	_, err := r.db.ExecContext(ctx, query,
		emobile.Name,
		emobile.PowerPointID,
		emobile.BatteryID,
		emobile.ChargerSystemID,
		emobile.ChassisID,
		emobile.BodyID,
		emobile.ElectronicsID,
		emobile.ID,
	)
	return err
}

func (r *EmobilePostgres) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM emobile WHERE emobile_id = $1`, id)
	return err
}

func (r *EmobilePostgres) List(ctx context.Context) ([]*domain.Emobile, error) {
	var rows []domain.Emobile
	query := `SELECT emobile_id, emobile_name, power_point_id, battery_id,
       charger_system_id, chassis_id, body_id, electronics_id FROM emobile`
	err := r.db.SelectContext(ctx, &rows, query)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.Emobile, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}
	return result, nil
}
