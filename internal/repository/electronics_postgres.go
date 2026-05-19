package repository

import (
	"context"

	"MISPRIS/internal/domain"

	"github.com/jmoiron/sqlx"
)

type ElectronicsPostgres struct {
	db *sqlx.DB
}

func NewElectronicsPostgres(db *sqlx.DB) *ElectronicsPostgres {
	return &ElectronicsPostgres{db: db}
}

func (r *ElectronicsPostgres) GetByID(ctx context.Context, id string) (*domain.Electronics, error) {
	var e domain.Electronics
	if err := r.db.GetContext(ctx, &e,
		`SELECT * FROM electronics WHERE electronics_id = $1`, id,
	); err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *ElectronicsPostgres) CreateTx(ctx context.Context, tx *sqlx.Tx, e *domain.Electronics) (string, error) {
	var electronicsID string
	err := tx.QueryRowContext(ctx,
		`INSERT INTO electronics (controller_id, sensor_id, wiring_id)
         VALUES ($1, $2, $3) RETURNING electronics_id`,
		e.ControllerID, e.SensorID, e.WiringID,
	).Scan(&electronicsID)
	return electronicsID, err
}

func (r *ElectronicsPostgres) Update(ctx context.Context, e *domain.Electronics) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE electronics SET controller_id=$1, sensor_id=$2, wiring_id=$3 WHERE electronics_id=$4`,
		e.ControllerID, e.SensorID, e.WiringID, e.ID,
	)
	return err
}

func (r *ElectronicsPostgres) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM electronics WHERE electronics_id = $1`, id,
	)
	return err
}

func (r *ElectronicsPostgres) List(ctx context.Context) ([]*domain.Electronics, error) {
	var rows []domain.Electronics
	if err := r.db.SelectContext(ctx, &rows, `SELECT * FROM electronics`); err != nil {
		return nil, err
	}
	result := make([]*domain.Electronics, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}
	return result, nil
}

//controllers

type ControllerPostgres struct{ db *sqlx.DB }

func NewControllerPostgres(db *sqlx.DB) *ControllerPostgres {
	return &ControllerPostgres{db: db}
}

func (r *ControllerPostgres) GetByID(ctx context.Context, id string) (*domain.Controller, error) {
	var c domain.Controller
	if err := r.db.GetContext(ctx, &c,
		`SELECT * FROM controllers WHERE controller_id = $1`, id,
	); err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ControllerPostgres) Create(ctx context.Context, c *domain.Controller) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO controllers (controller_name, controller_info)
         VALUES ($1, $2) RETURNING controller_id`,
		c.Name, c.Info,
	).Scan(&id)
	return id, err
}

func (r *ControllerPostgres) Update(ctx context.Context, c *domain.Controller) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE controllers SET controller_name=$1, controller_info=$2 WHERE controller_id=$3`,
		c.Name, c.Info, c.ID,
	)
	return err
}

func (r *ControllerPostgres) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM controllers WHERE controller_id = $1`, id,
	)
	return err
}

func (r *ControllerPostgres) List(ctx context.Context) ([]*domain.Controller, error) {
	var rows []domain.Controller
	if err := r.db.SelectContext(ctx, &rows, `SELECT * FROM controllers`); err != nil {
		return nil, err
	}
	result := make([]*domain.Controller, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}
	return result, nil
}

//sensors

type SensorPostgres struct{ db *sqlx.DB }

func NewSensorPostgres(db *sqlx.DB) *SensorPostgres {
	return &SensorPostgres{db: db}
}

func (r *SensorPostgres) GetByID(ctx context.Context, id string) (*domain.Sensor, error) {
	var s domain.Sensor
	if err := r.db.GetContext(ctx, &s,
		`SELECT * FROM sensors WHERE sensor_id = $1`, id,
	); err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SensorPostgres) Create(ctx context.Context, s *domain.Sensor) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO sensors (sensor_name, sensor_info)
         VALUES ($1, $2) RETURNING sensor_id`,
		s.Name, s.Info,
	).Scan(&id)
	return id, err
}

func (r *SensorPostgres) Update(ctx context.Context, s *domain.Sensor) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE sensors SET sensor_name=$1, sensor_info=$2 WHERE sensor_id=$3`,
		s.Name, s.Info, s.ID,
	)
	return err
}

func (r *SensorPostgres) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM sensors WHERE sensor_id = $1`, id,
	)
	return err
}

func (r *SensorPostgres) List(ctx context.Context) ([]*domain.Sensor, error) {
	var rows []domain.Sensor
	if err := r.db.SelectContext(ctx, &rows, `SELECT * FROM sensors`); err != nil {
		return nil, err
	}
	result := make([]*domain.Sensor, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}
	return result, nil
}

//wiring

type WiringPostgres struct{ db *sqlx.DB }

func NewWiringPostgres(db *sqlx.DB) *WiringPostgres {
	return &WiringPostgres{db: db}
}

func (r *WiringPostgres) GetByID(ctx context.Context, id string) (*domain.Wiring, error) {
	var w domain.Wiring
	if err := r.db.GetContext(ctx, &w,
		`SELECT * FROM wiring WHERE wiring_id = $1`, id,
	); err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *WiringPostgres) Create(ctx context.Context, w *domain.Wiring) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO wiring (wiring_name, wiring_info)
         VALUES ($1, $2) RETURNING wiring_id`,
		w.Name, w.Info,
	).Scan(&id)
	return id, err
}

func (r *WiringPostgres) Update(ctx context.Context, w *domain.Wiring) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE wiring SET wiring_name=$1, wiring_info=$2 WHERE wiring_id=$3`,
		w.Name, w.Info, w.ID,
	)
	return err
}

func (r *WiringPostgres) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM wiring WHERE wiring_id = $1`, id,
	)
	return err
}

func (r *WiringPostgres) List(ctx context.Context) ([]*domain.Wiring, error) {
	var rows []domain.Wiring
	if err := r.db.SelectContext(ctx, &rows, `SELECT * FROM wiring`); err != nil {
		return nil, err
	}
	result := make([]*domain.Wiring, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}
	return result, nil
}
