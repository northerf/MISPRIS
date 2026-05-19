package repository

import (
	"context"

	"MISPRIS/internal/domain"

	"github.com/jmoiron/sqlx"
)

// ChassisPostgres
type ChassisPostgres struct {
	db *sqlx.DB
}

func NewChassisPostgres(db *sqlx.DB) *ChassisPostgres {
	return &ChassisPostgres{db: db}
}

func (r *ChassisPostgres) GetByID(ctx context.Context, id string) (*domain.Chassis, error) {
	var c domain.Chassis
	err := r.db.GetContext(ctx, &c,
		`SELECT * FROM chassis WHERE chassis_id = $1`, id,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ChassisPostgres) CreateTx(ctx context.Context, tx *sqlx.Tx, c *domain.Chassis) (string, error) {
	var chassisID string
	err := tx.QueryRowContext(ctx,
		`INSERT INTO chassis (frame_id, suspension_id, break_system_id)
		 VALUES ($1, $2, $3) RETURNING chassis_id`,
		c.FrameID, c.SuspensionID, c.BreakSystemID,
	).Scan(&chassisID)
	return chassisID, err
}

func (r *ChassisPostgres) Update(ctx context.Context, c *domain.Chassis) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE chassis SET frame_id=$1, suspension_id=$2, break_system_id=$3 WHERE chassis_id=$4`,
		c.FrameID, c.SuspensionID, c.BreakSystemID, c.ID,
	)
	return err
}

func (r *ChassisPostgres) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM chassis WHERE chassis_id = $1`, id,
	)
	return err
}

func (r *ChassisPostgres) List(ctx context.Context) ([]*domain.Chassis, error) {
	var rows []domain.Chassis
	err := r.db.SelectContext(ctx, &rows, `SELECT * FROM chassis`)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.Chassis, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}
	return result, nil
}

// FramePostgres
type FramePostgres struct {
	db *sqlx.DB
}

func NewFramePostgres(db *sqlx.DB) *FramePostgres {
	return &FramePostgres{db: db}
}

func (r *FramePostgres) GetByID(ctx context.Context, id string) (*domain.Frame, error) {
	var f domain.Frame
	err := r.db.GetContext(ctx, &f,
		`SELECT * FROM frame WHERE frame_id = $1`, id,
	)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *FramePostgres) Create(ctx context.Context, f *domain.Frame) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO frame (frame_name, frame_info)
		 VALUES ($1, $2) RETURNING frame_id`,
		f.Name, f.Info,
	).Scan(&id)
	return id, err
}

func (r *FramePostgres) Update(ctx context.Context, f *domain.Frame) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE frame SET frame_name=$1, frame_info=$2 WHERE frame_id=$3`,
		f.Name, f.Info, f.ID,
	)
	return err
}

func (r *FramePostgres) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM frame WHERE frame_id = $1`, id,
	)
	return err
}

func (r *FramePostgres) List(ctx context.Context) ([]*domain.Frame, error) {
	var rows []domain.Frame
	err := r.db.SelectContext(ctx, &rows, `SELECT * FROM frame`)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.Frame, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}
	return result, nil
}

// SuspensionPostgres
type SuspensionPostgres struct {
	db *sqlx.DB
}

func NewSuspensionPostgres(db *sqlx.DB) *SuspensionPostgres {
	return &SuspensionPostgres{db: db}
}

func (r *SuspensionPostgres) GetByID(ctx context.Context, id string) (*domain.Suspension, error) {
	var s domain.Suspension
	err := r.db.GetContext(ctx, &s,
		`SELECT * FROM suspension WHERE suspension_id = $1`, id,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SuspensionPostgres) Create(ctx context.Context, s *domain.Suspension) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO suspension (suspension_name, suspension_info)
		 VALUES ($1, $2) RETURNING suspension_id`,
		s.Name, s.Info,
	).Scan(&id)
	return id, err
}

func (r *SuspensionPostgres) Update(ctx context.Context, s *domain.Suspension) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE suspension SET suspension_name=$1, suspension_info=$2 WHERE suspension_id=$3`,
		s.Name, s.Info, s.ID,
	)
	return err
}

func (r *SuspensionPostgres) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM suspension WHERE suspension_id = $1`, id,
	)
	return err
}

func (r *SuspensionPostgres) List(ctx context.Context) ([]*domain.Suspension, error) {
	var rows []domain.Suspension
	err := r.db.SelectContext(ctx, &rows, `SELECT * FROM suspension`)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.Suspension, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}
	return result, nil
}

// BreakSystemPostgres
type BreakSystemPostgres struct {
	db *sqlx.DB
}

func NewBreakSystemPostgres(db *sqlx.DB) *BreakSystemPostgres {
	return &BreakSystemPostgres{db: db}
}

func (r *BreakSystemPostgres) GetByID(ctx context.Context, id string) (*domain.BreakSystem, error) {
	var b domain.BreakSystem
	err := r.db.GetContext(ctx, &b,
		`SELECT * FROM break_system WHERE break_system_id = $1`, id,
	)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BreakSystemPostgres) Create(ctx context.Context, b *domain.BreakSystem) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO break_system (break_system_name, break_info)
		 VALUES ($1, $2) RETURNING break_system_id`,
		b.Name, b.Info,
	).Scan(&id)
	return id, err
}

func (r *BreakSystemPostgres) Update(ctx context.Context, b *domain.BreakSystem) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE break_system SET break_system_name=$1, break_info=$2 WHERE break_system_id=$3`,
		b.Name, b.Info, b.ID,
	)
	return err
}

func (r *BreakSystemPostgres) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM break_system WHERE break_system_id = $1`, id,
	)
	return err
}

func (r *BreakSystemPostgres) List(ctx context.Context) ([]*domain.BreakSystem, error) {
	var rows []domain.BreakSystem
	err := r.db.SelectContext(ctx, &rows, `SELECT * FROM break_system`)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.BreakSystem, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}
	return result, nil
}
