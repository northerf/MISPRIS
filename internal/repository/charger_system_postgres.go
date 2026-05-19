package repository

import (
	"context"

	"MISPRIS/internal/domain"

	"github.com/jmoiron/sqlx"
)

// ChargerSystemPostgres
type ChargerSystemPostgres struct {
	db *sqlx.DB
}

func NewChargerSystemPostgres(db *sqlx.DB) *ChargerSystemPostgres {
	return &ChargerSystemPostgres{db: db}
}

func (r *ChargerSystemPostgres) GetByID(ctx context.Context, id string) (*domain.ChargerSystem, error) {
	var cs domain.ChargerSystem
	err := r.db.GetContext(ctx, &cs,
		`SELECT * FROM charger_system WHERE charger_system_id = $1`, id,
	)
	if err != nil {
		return nil, err
	}
	return &cs, nil
}

func (r *ChargerSystemPostgres) CreateTx(ctx context.Context, tx *sqlx.Tx, cs *domain.ChargerSystem) (string, error) {
	var charger_systemID string
	err := tx.QueryRowContext(ctx,
		`INSERT INTO charger_system (charger_id, connector_id)
		 VALUES ($1, $2) RETURNING charger_system_id`,
		cs.ChargerID, cs.ConnectorID,
	).Scan(&charger_systemID)
	return charger_systemID, err
}

func (r *ChargerSystemPostgres) Update(ctx context.Context, cs *domain.ChargerSystem) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE charger_system SET charger_id=$1, connector_id=$2 WHERE charger_system_id=$3`,
		cs.ChargerID, cs.ConnectorID, cs.ID,
	)
	return err
}

func (r *ChargerSystemPostgres) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM charger_system WHERE charger_system_id = $1`, id,
	)
	return err
}

func (r *ChargerSystemPostgres) List(ctx context.Context) ([]*domain.ChargerSystem, error) {
	var rows []domain.ChargerSystem
	err := r.db.SelectContext(ctx, &rows, `SELECT * FROM charger_system`)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.ChargerSystem, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}
	return result, nil
}

// ChargerPostgres
type ChargerPostgres struct {
	db *sqlx.DB
}

func NewChargerPostgres(db *sqlx.DB) *ChargerPostgres {
	return &ChargerPostgres{db: db}
}

func (r *ChargerPostgres) GetByID(ctx context.Context, id string) (*domain.Charger, error) {
	var c domain.Charger
	err := r.db.GetContext(ctx, &c,
		`SELECT * FROM charger WHERE charger_id = $1`, id,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *ChargerPostgres) Create(ctx context.Context, c *domain.Charger) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO charger (charger_name, charger_info)
		 VALUES ($1, $2) RETURNING charger_id`,
		c.Name, c.Info,
	).Scan(&id)
	return id, err
}

func (r *ChargerPostgres) Update(ctx context.Context, c *domain.Charger) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE charger SET charger_name=$1, charger_info=$2 WHERE charger_id=$3`,
		c.Name, c.Info, c.ID,
	)
	return err
}

func (r *ChargerPostgres) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM charger WHERE charger_id = $1`, id,
	)
	return err
}

func (r *ChargerPostgres) List(ctx context.Context) ([]*domain.Charger, error) {
	var rows []domain.Charger
	err := r.db.SelectContext(ctx, &rows, `SELECT * FROM charger`)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.Charger, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}
	return result, nil
}

// ConnectorPostgres
type ConnectorPostgres struct {
	db *sqlx.DB
}

func NewConnectorPostgres(db *sqlx.DB) *ConnectorPostgres {
	return &ConnectorPostgres{db: db}
}

func (r *ConnectorPostgres) GetByID(ctx context.Context, id string) (*domain.Connector, error) {
	var conn domain.Connector
	err := r.db.GetContext(ctx, &conn,
		`SELECT * FROM connector WHERE connector_id = $1`, id,
	)
	if err != nil {
		return nil, err
	}
	return &conn, nil
}

func (r *ConnectorPostgres) Create(ctx context.Context, conn *domain.Connector) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO connector (connector_name, connector_info)
		 VALUES ($1, $2) RETURNING connector_id`,
		conn.Name, conn.Info,
	).Scan(&id)
	return id, err
}

func (r *ConnectorPostgres) Update(ctx context.Context, conn *domain.Connector) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE connector SET connector_name=$1, connector_info=$2 WHERE connector_id=$3`,
		conn.Name, conn.Info, conn.ID,
	)
	return err
}

func (r *ConnectorPostgres) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM connector WHERE connector_id = $1`, id,
	)
	return err
}

func (r *ConnectorPostgres) List(ctx context.Context) ([]*domain.Connector, error) {
	var rows []domain.Connector
	err := r.db.SelectContext(ctx, &rows, `SELECT * FROM connector`)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.Connector, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}
	return result, nil
}
