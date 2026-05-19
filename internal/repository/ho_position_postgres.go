package repository

import (
	"context"

	"MISPRIS/internal/domain"

	"github.com/jmoiron/sqlx"
)

type HoPositionPostgres struct {
	db *sqlx.DB
}

func NewHoPositionPostgres(db *sqlx.DB) *HoPositionPostgres {
	return &HoPositionPostgres{db: db}
}

// ListByHo возвращает позиции ХО с присоединённым названием ТС.
// Если hoID не пустой — фильтрует по экземпляру ХО.
func (r *HoPositionPostgres) ListByHo(ctx context.Context, hoID string) ([]*domain.HoPositionFull, error) {
	var items []*domain.HoPositionFull
	var err error
	baseSQL := `SELECT p.id, p.ho_id, p.emobile_id,
	              COALESCE(e.name,'') AS emobile_name,
	              p.quantity,
	              COALESCE(p.unit_price,0) AS unit_price,
	              COALESCE(p.total_price,0) AS total_price,
	              COALESCE(p.note,'') AS note,
	              p.position_num
	             FROM ho_position p
	             LEFT JOIN emobile e ON e.emobile_id = p.emobile_id`
	if hoID == "" {
		err = r.db.SelectContext(ctx, &items, baseSQL+` ORDER BY p.ho_id, p.position_num`)
	} else {
		err = r.db.SelectContext(ctx, &items,
			baseSQL+` WHERE p.ho_id = $1 ORDER BY p.position_num`, hoID)
	}
	return items, err
}

// Create вызывает SQL-процедуру add_ho_position, которая автоматически пересчитывает итоговую сумму.
func (r *HoPositionPostgres) Create(ctx context.Context, p *domain.HoPosition) (string, error) {
	_, err := r.db.ExecContext(ctx,
		`CALL add_ho_position($1, $2, $3, NULLIF($4,0))`,
		p.HoID, p.EmobileID, p.Quantity, p.UnitPrice)
	if err != nil {
		return "", err
	}
	// Возвращаем id последней добавленной позиции в данном ХО
	var id string
	err = r.db.QueryRowContext(ctx,
		`SELECT id FROM ho_position WHERE ho_id=$1 AND emobile_id=$2 ORDER BY id DESC LIMIT 1`,
		p.HoID, p.EmobileID,
	).Scan(&id)
	return id, err
}

// Update выполняет прямое обновление позиции (без пересчёта итога).
func (r *HoPositionPostgres) Update(ctx context.Context, p *domain.HoPosition) error {
	query := `UPDATE ho_position
	          SET quantity=$1, unit_price=NULLIF($2,0), note=NULLIF($3,''),
	              total_price = quantity * COALESCE(NULLIF($2,0), unit_price)
	          WHERE id = $4`
	_, err := r.db.ExecContext(ctx, query, p.Quantity, p.UnitPrice, p.Note, p.ID)
	return err
}

func (r *HoPositionPostgres) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM ho_position WHERE id = $1`, id)
	return err
}
