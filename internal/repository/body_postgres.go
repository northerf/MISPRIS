package repository

import (
	"context"

	"MISPRIS/internal/domain"

	"github.com/jmoiron/sqlx"
)

type BodyPostgres struct {
	db *sqlx.DB
}

func NewBodyPostgres(db *sqlx.DB) *BodyPostgres {
	return &BodyPostgres{db: db}
}

func (r *BodyPostgres) GetByID(ctx context.Context, id string) (*domain.Body, error) {
	var b domain.Body
	if err := r.db.GetContext(ctx, &b,
		`SELECT * FROM body WHERE body_id = $1`, id,
	); err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BodyPostgres) CreateTx(ctx context.Context, tx *sqlx.Tx, b *domain.Body) (string, error) {
	var bodyID string
	err := tx.QueryRowContext(ctx,
		`INSERT INTO body (carcass_id, doors_id, wings_id)
         VALUES ($1, $2, $3) RETURNING body_id`,
		b.CarcassID, b.DoorsID, b.WingsID,
	).Scan(&bodyID)
	return bodyID, err
}

func (r *BodyPostgres) Update(ctx context.Context, b *domain.Body) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE body SET carcass_id=$1, doors_id=$2, wings_id=$3 WHERE body_id=$4`,
		b.CarcassID, b.DoorsID, b.WingsID, b.ID,
	)
	return err
}

func (r *BodyPostgres) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM body WHERE body_id = $1`, id,
	)
	return err
}

func (r *BodyPostgres) List(ctx context.Context) ([]*domain.Body, error) {
	var rows []domain.Body
	if err := r.db.SelectContext(ctx, &rows, `SELECT * FROM body`); err != nil {
		return nil, err
	}
	result := make([]*domain.Body, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}
	return result, nil

}

//carcass

type CarcassPostgres struct {
	db *sqlx.DB
}

func NewCarcassPostgres(db *sqlx.DB) *CarcassPostgres {
	return &CarcassPostgres{db: db}
}

func (r *CarcassPostgres) GetByID(ctx context.Context, id string) (*domain.Carcass, error) {
	var c domain.Carcass
	if err := r.db.GetContext(ctx, &c,
		`SELECT * FROM carcass WHERE carcass_id = $1`, id,
	); err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CarcassPostgres) Create(ctx context.Context, c *domain.Carcass) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO carcass (carcass_name, carcass_info)
         VALUES ($1, $2) RETURNING carcass_id`,
		c.Name, c.Info,
	).Scan(&id)
	return id, err
}

func (r *CarcassPostgres) Update(ctx context.Context, c *domain.Carcass) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE carcass SET carcass_name=$1, carcass_info=$2 WHERE carcass_id=$3`,
		c.Name, c.Info, c.ID,
	)
	return err
}

func (r *CarcassPostgres) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM carcass WHERE carcass_id = $1`, id,
	)
	return err
}

func (r *CarcassPostgres) List(ctx context.Context) ([]*domain.Carcass, error) {
	var rows []domain.Carcass
	if err := r.db.SelectContext(ctx, &rows, `SELECT * FROM carcass`); err != nil {
		return nil, err
	}
	result := make([]*domain.Carcass, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}
	return result, nil
}

//doors

type DoorsPostgres struct {
	db *sqlx.DB
}

func NewDoorsPostgres(db *sqlx.DB) *DoorsPostgres {
	return &DoorsPostgres{db: db}
}

func (r *DoorsPostgres) GetByID(ctx context.Context, id string) (*domain.Doors, error) {
	var d domain.Doors
	if err := r.db.GetContext(ctx, &d,
		`SELECT * FROM doors WHERE doors_id = $1`, id,
	); err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DoorsPostgres) Create(ctx context.Context, d *domain.Doors) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO doors (doors_name, doors_info)
         VALUES ($1, $2) RETURNING doors_id`,
		d.Name, d.Info,
	).Scan(&id)
	return id, err
}

func (r *DoorsPostgres) Update(ctx context.Context, d *domain.Doors) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE doors SET doors_name=$1, doors_info=$2 WHERE doors_id=$3`,
		d.Name, d.Info, d.ID,
	)
	return err
}

func (r *DoorsPostgres) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM doors WHERE doors_id = $1`, id,
	)
	return err
}

func (r *DoorsPostgres) List(ctx context.Context) ([]*domain.Doors, error) {
	var rows []domain.Doors
	if err := r.db.SelectContext(ctx, &rows, `SELECT * FROM doors`); err != nil {
		return nil, err
	}
	result := make([]*domain.Doors, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}
	return result, nil
}

// wings
type WingsPostgres struct {
	db *sqlx.DB
}

func NewWingsPostgres(db *sqlx.DB) *WingsPostgres {
	return &WingsPostgres{db: db}
}

func (r *WingsPostgres) GetByID(ctx context.Context, id string) (*domain.Wings, error) {
	var w domain.Wings
	if err := r.db.GetContext(ctx, &w,
		`SELECT * FROM wings WHERE wings_id = $1`, id,
	); err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *WingsPostgres) Create(ctx context.Context, w *domain.Wings) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO wings (wings_name, wings_info)
         VALUES ($1, $2) RETURNING wings_id`,
		w.Name, w.Info,
	).Scan(&id)
	return id, err
}

func (r *WingsPostgres) Update(ctx context.Context, w *domain.Wings) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE wings SET wings_name=$1, wings_info=$2 WHERE wings_id=$3`,
		w.Name, w.Info, w.ID,
	)
	return err
}

func (r *WingsPostgres) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM wings WHERE wings_id = $1`, id,
	)
	return err
}

func (r *WingsPostgres) List(ctx context.Context) ([]*domain.Wings, error) {
	var rows []domain.Wings
	if err := r.db.SelectContext(ctx, &rows, `SELECT * FROM wings`); err != nil {
		return nil, err
	}
	result := make([]*domain.Wings, len(rows))
	for i := range rows {
		result[i] = &rows[i]
	}
	return result, nil
}
