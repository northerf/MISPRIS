package repository

import (
	"context"
	"time"

	"MISPRIS/internal/domain"

	"github.com/jmoiron/sqlx"
)

type EnumClassPostgres struct {
	db *sqlx.DB
}

func NewEnumClassPostgres(db *sqlx.DB) *EnumClassPostgres {
	return &EnumClassPostgres{db: db}
}

func (r *EnumClassPostgres) List(ctx context.Context) ([]*domain.EnumClass, error) {
	var items []*domain.EnumClass
	query := `SELECT enum_class_id, name, COALESCE(component_type,'') AS component_type,
	           created_at, updated_at
	          FROM enum_class ORDER BY enum_class_id`
	err := r.db.SelectContext(ctx, &items, query)
	return items, err
}

func (r *EnumClassPostgres) GetByID(ctx context.Context, id string) (*domain.EnumClass, error) {
	var item domain.EnumClass
	query := `SELECT enum_class_id, name, COALESCE(component_type,'') AS component_type,
	           created_at, updated_at
	          FROM enum_class WHERE enum_class_id = $1`
	err := r.db.GetContext(ctx, &item, query, id)
	return &item, err
}

func (r *EnumClassPostgres) Create(ctx context.Context, ec *domain.EnumClass) (string, error) {
	var id string
	query := `INSERT INTO enum_class (name, component_type, created_at, updated_at)
	          VALUES ($1, NULLIF($2,''), NOW(), NOW()) RETURNING enum_class_id`
	err := r.db.QueryRowContext(ctx, query, ec.Name, ec.ComponentType).Scan(&id)
	return id, err
}

func (r *EnumClassPostgres) Update(ctx context.Context, ec *domain.EnumClass) error {
	query := `UPDATE enum_class SET name=$1, component_type=NULLIF($2,''), updated_at=$3
	          WHERE enum_class_id = $4`
	_, err := r.db.ExecContext(ctx, query, ec.Name, ec.ComponentType, time.Now(), ec.ID)
	return err
}

func (r *EnumClassPostgres) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM enum_class WHERE enum_class_id = $1`, id)
	return err
}

// GetValues — прямой запрос вместо SQL-функции (обходим несовпадение колонок).
// Возвращает позиции упорядоченные по order_num.
func (r *EnumClassPostgres) GetValues(ctx context.Context, id string) ([]*domain.EnumPosition, error) {
	var positions []*domain.EnumPosition
	query := `SELECT enum_position_id, enum_class_id, value, order_num,
	           created_at, updated_at
	          FROM enum_position WHERE enum_class_id = $1 ORDER BY order_num`
	err := r.db.SelectContext(ctx, &positions, query, id)
	return positions, err
}

// ValidateValue вызывает SQL-функцию validate_enum_value.
func (r *EnumClassPostgres) ValidateValue(ctx context.Context, enumClassID, value string) (bool, error) {
	var result bool
	err := r.db.QueryRowContext(ctx,
		`SELECT validate_enum_value($1, $2)`, enumClassID, value,
	).Scan(&result)
	return result, err
}
