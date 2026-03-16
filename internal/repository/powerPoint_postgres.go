package repository

import (
	"context"

	"MISPRIS/internal/domain"

	"github.com/jmoiron/sqlx"
)

// ---- PowerPoint ----

type PowerPointPostgres struct {
	db *sqlx.DB
}

func NewPowerPointPostgres(db *sqlx.DB) *PowerPointPostgres {
	return &PowerPointPostgres{db: db}
}

func (r *PowerPointPostgres) GetByID(ctx context.Context, id int64) (*domain.PowerPoint, error) {
	var p domain.PowerPoint
	query := `SELECT power_point_id, engine_id, inverter_id, gearbox_id
	          FROM power_point WHERE power_point_id = $1`
	err := r.db.GetContext(ctx, &p, query, id)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PowerPointPostgres) Create(ctx context.Context, p *domain.PowerPoint) (int64, error) {
	var id int64
	query := `INSERT INTO power_point (engine_id, inverter_id, gearbox_id)
	          VALUES ($1, $2, $3) RETURNING power_point_id`
	err := r.db.QueryRowContext(ctx, query, p.EngineID, p.InverterID, p.GearboxID).Scan(&id)
	return id, err
}

func (r *PowerPointPostgres) Update(ctx context.Context, p *domain.PowerPoint) error {
	query := `UPDATE power_point
	          SET engine_id = $1, inverter_id = $2, gearbox_id = $3
	          WHERE power_point_id = $4`
	_, err := r.db.ExecContext(ctx, query, p.EngineID, p.InverterID, p.GearboxID, p.ID)
	return err
}

func (r *PowerPointPostgres) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM power_point WHERE power_point_id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PowerPointPostgres) List(ctx context.Context) ([]*domain.PowerPoint, error) {
	var rows []domain.PowerPoint
	query := `SELECT power_point_id, engine_id, inverter_id, gearbox_id FROM power_point`
	err := r.db.SelectContext(ctx, &rows, query)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.PowerPoint, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}
	return result, nil
}

// ---- Engine ----

type EnginePostgres struct {
	db *sqlx.DB
}

func NewEnginePostgres(db *sqlx.DB) *EnginePostgres {
	return &EnginePostgres{db: db}
}

func (r *EnginePostgres) GetByID(ctx context.Context, id int64) (*domain.Engine, error) {
	var e domain.Engine
	query := `SELECT engine_id, engine_name, engine_type, engine_info FROM engine WHERE engine_id = $1`
	err := r.db.GetContext(ctx, &e, query, id)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func (r *EnginePostgres) Create(ctx context.Context, e *domain.Engine) (int64, error) {
	var id int64
	query := `INSERT INTO engine (engine_name, engine_type, engine_info) VALUES ($1, $2, $3) RETURNING engine_id`
	err := r.db.QueryRowContext(ctx, query, e.EngineName, e.EngineType, e.EngineInfo).Scan(&id)
	return id, err
}

func (r *EnginePostgres) Update(ctx context.Context, e *domain.Engine) error {
	query := `UPDATE engine SET engine_name = $1, engine_type = $2, engine_info = $3 WHERE engine_id = $4`
	_, err := r.db.ExecContext(ctx, query, e.EngineName, e.EngineType, e.EngineInfo, e.ID)
	return err
}

func (r *EnginePostgres) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM engine WHERE engine_id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *EnginePostgres) List(ctx context.Context) ([]*domain.Engine, error) {
	var rows []domain.Engine
	query := `SELECT engine_id, engine_name, engine_type, engine_info FROM engine`
	err := r.db.SelectContext(ctx, &rows, query)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.Engine, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}
	return result, nil
}

// ---- Inverter ----

type InverterPostgres struct {
	db *sqlx.DB
}

func NewInverterPostgres(db *sqlx.DB) *InverterPostgres {
	return &InverterPostgres{db: db}
}

func (r *InverterPostgres) GetByID(ctx context.Context, id int64) (*domain.Inverter, error) {
	var i domain.Inverter
	query := `SELECT inverter_id, inverter_name, inverter_info FROM inverter WHERE inverter_id = $1`
	err := r.db.GetContext(ctx, &i, query, id)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func (r *InverterPostgres) Create(ctx context.Context, i *domain.Inverter) (int64, error) {
	var id int64
	query := `INSERT INTO inverter (inverter_name, inverter_info) VALUES ($1, $2) RETURNING inverter_id`
	err := r.db.QueryRowContext(ctx, query, i.InverterName, i.InverterInfo).Scan(&id)
	return id, err
}

func (r *InverterPostgres) Update(ctx context.Context, i *domain.Inverter) error {
	query := `UPDATE inverter SET inverter_name = $1, inverter_info = $2 WHERE inverter_id = $3`
	_, err := r.db.ExecContext(ctx, query, i.InverterName, i.InverterInfo, i.ID)
	return err
}

func (r *InverterPostgres) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM inverter WHERE inverter_id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *InverterPostgres) List(ctx context.Context) ([]*domain.Inverter, error) {
	var rows []domain.Inverter
	query := `SELECT inverter_id, inverter_name, inverter_info FROM inverter`
	err := r.db.SelectContext(ctx, &rows, query)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.Inverter, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}
	return result, nil
}

// ---- Gearbox ----

type GearboxPostgres struct {
	db *sqlx.DB
}

func NewGearboxPostgres(db *sqlx.DB) *GearboxPostgres {
	return &GearboxPostgres{db: db}
}

func (r *GearboxPostgres) GetByID(ctx context.Context, id int64) (*domain.Gearbox, error) {
	var g domain.Gearbox
	query := `SELECT gearbox_id, gearbox_name, gearbox_info FROM gearbox WHERE gearbox_id = $1`
	err := r.db.GetContext(ctx, &g, query, id)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *GearboxPostgres) Create(ctx context.Context, g *domain.Gearbox) (int64, error) {
	var id int64
	query := `INSERT INTO gearbox (gearbox_name, gearbox_info) VALUES ($1, $2) RETURNING gearbox_id`
	err := r.db.QueryRowContext(ctx, query, g.GearboxName, g.GearboxInfo).Scan(&id)
	return id, err
}

func (r *GearboxPostgres) Update(ctx context.Context, g *domain.Gearbox) error {
	query := `UPDATE gearbox SET gearbox_name = $1, gearbox_info = $2 WHERE gearbox_id = $3`
	_, err := r.db.ExecContext(ctx, query, g.GearboxName, g.GearboxInfo, g.ID)
	return err
}

func (r *GearboxPostgres) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM gearbox WHERE gearbox_id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *GearboxPostgres) List(ctx context.Context) ([]*domain.Gearbox, error) {
	var rows []domain.Gearbox
	query := `SELECT gearbox_id, gearbox_name, gearbox_info FROM gearbox`
	err := r.db.SelectContext(ctx, &rows, query)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.Gearbox, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}
	return result, nil
}
