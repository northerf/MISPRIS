package repository

import (
	"context"

	"MISPRIS/internal/domain"

	"github.com/jmoiron/sqlx"
)

type HoActorPostgres struct {
	db *sqlx.DB
}

func NewHoActorPostgres(db *sqlx.DB) *HoActorPostgres {
	return &HoActorPostgres{db: db}
}

// ListByHo возвращает всех акторов; если hoID не пустой — фильтрует по экземпляру ХО.
func (r *HoActorPostgres) ListByHo(ctx context.Context, hoID string) ([]*domain.HoActor, error) {
	var items []*domain.HoActor
	var err error
	if hoID == "" {
		err = r.db.SelectContext(ctx, &items,
			`SELECT id, ho_id, ho_role_id, shd_id FROM ho_actor ORDER BY id`)
	} else {
		err = r.db.SelectContext(ctx, &items,
			`SELECT id, ho_id, ho_role_id, shd_id FROM ho_actor WHERE ho_id = $1 ORDER BY id`, hoID)
	}
	return items, err
}

// Create вызывает SQL-процедуру set_ho_actor, которая валидирует и назначает актора.
func (r *HoActorPostgres) Create(ctx context.Context, ha *domain.HoActor) (string, error) {
	_, err := r.db.ExecContext(ctx, `CALL set_ho_actor($1, $2, $3)`,
		ha.HoID, ha.HoRoleID, ha.ShdID)
	if err != nil {
		return "", err
	}
	// Получаем только что созданную запись
	var id string
	err = r.db.QueryRowContext(ctx,
		`SELECT id FROM ho_actor WHERE ho_id=$1 AND ho_role_id=$2`,
		ha.HoID, ha.HoRoleID,
	).Scan(&id)
	return id, err
}

func (r *HoActorPostgres) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM ho_actor WHERE id = $1`, id)
	return err
}
