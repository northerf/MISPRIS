package repository

import (
	"context"

	"MISPRIS/internal/domain"

	"github.com/jmoiron/sqlx"
)

type HoParameterValuePostgres struct {
	db *sqlx.DB
}

func NewHoParameterValuePostgres(db *sqlx.DB) *HoParameterValuePostgres {
	return &HoParameterValuePostgres{db: db}
}

// ListByHo возвращает значения параметров с именем и типом из таблицы parameter.
// Если hoID не пустой — фильтрует по экземпляру ХО.
func (r *HoParameterValuePostgres) ListByHo(ctx context.Context, hoID string) ([]*domain.HoParameterValueFull, error) {
	var items []*domain.HoParameterValueFull
	base := `
		SELECT
			hpv.id,
			hpv.ho_id,
			hpv.ho_class_parameter_id,
			COALESCE(p.name, '')                AS param_name,
			COALESCE(p.param_type, '')          AS param_type,
			COALESCE(p.measuring_unit, '')      AS measuring_unit,
			COALESCE(hpv.val_real, 0)          AS val_real,
			COALESCE(hpv.val_int, 0)           AS val_int,
			COALESCE(hpv.val_str, '')          AS val_str,
			COALESCE(hpv.val_date::text, '')   AS val_date,
			COALESCE(hpv.enum_val_id::text,'') AS enum_val_id
		FROM ho_parameter_value hpv
		LEFT JOIN ho_class_parameter hcp ON hpv.ho_class_parameter_id = hcp.id
		LEFT JOIN parameter p ON hcp.parameter_id = p.parameter_id`
	var err error
	if hoID == "" {
		err = r.db.SelectContext(ctx, &items, base+` ORDER BY hpv.id`)
	} else {
		err = r.db.SelectContext(ctx, &items, base+` WHERE hpv.ho_id = $1 ORDER BY hpv.id`, hoID)
	}
	return items, err
}

func (r *HoParameterValuePostgres) GetByID(ctx context.Context, id string) (*domain.HoParameterValue, error) {
	var item domain.HoParameterValue
	query := `SELECT id, ho_id, ho_class_parameter_id,
	           COALESCE(val_real,0) AS val_real, COALESCE(val_int,0) AS val_int,
	           COALESCE(val_str,'') AS val_str, COALESCE(val_date::text,'') AS val_date,
	           COALESCE(enum_val_id::text,'') AS enum_val_id
	          FROM ho_parameter_value WHERE id = $1`
	err := r.db.GetContext(ctx, &item, query, id)
	return &item, err
}

// Create вызывает SQL-процедуру write_ho_par, которая валидирует и записывает значение.
func (r *HoParameterValuePostgres) Create(ctx context.Context, v *domain.HoParameterValue) (string, error) {
	_, err := r.db.ExecContext(ctx,
		`CALL write_ho_par($1, $2, NULLIF($3,0), NULLIF($4,0), NULLIF($5,''), NULLIF($6,'')::date, NULLIF($7,'')::int)`,
		v.HoID, v.HoClassParameterID, v.ValReal, v.ValInt, v.ValStr, v.ValDate, v.EnumValID)
	if err != nil {
		return "", err
	}
	var id string
	err = r.db.QueryRowContext(ctx,
		`SELECT id FROM ho_parameter_value WHERE ho_id=$1 AND ho_class_parameter_id=$2`,
		v.HoID, v.HoClassParameterID,
	).Scan(&id)
	return id, err
}

// Update вызывает SQL-процедуру write_ho_par повторно (UPSERT-семантика процедуры).
func (r *HoParameterValuePostgres) Update(ctx context.Context, v *domain.HoParameterValue) error {
	_, err := r.db.ExecContext(ctx,
		`CALL write_ho_par($1, $2, NULLIF($3,0), NULLIF($4,0), NULLIF($5,''), NULLIF($6,'')::date, NULLIF($7,'')::int)`,
		v.HoID, v.HoClassParameterID, v.ValReal, v.ValInt, v.ValStr, v.ValDate, v.EnumValID)
	return err
}

func (r *HoParameterValuePostgres) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM ho_parameter_value WHERE id = $1`, id)
	return err
}
