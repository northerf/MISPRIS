package repository

import (
	"context"

	"MISPRIS/internal/domain"

	"github.com/jmoiron/sqlx"
)

type HoClassDocumentPostgres struct {
	db *sqlx.DB
}

func NewHoClassDocumentPostgres(db *sqlx.DB) *HoClassDocumentPostgres {
	return &HoClassDocumentPostgres{db: db}
}

// ListByClass возвращает записи ho_class_document; если hoClassID не пустой — фильтрует по классу.
func (r *HoClassDocumentPostgres) ListByClass(ctx context.Context, hoClassID string) ([]*domain.HoClassDocument, error) {
	var items []*domain.HoClassDocument
	var err error
	if hoClassID == "" {
		err = r.db.SelectContext(ctx, &items,
			`SELECT id, ho_class_id, doc_class_id,
			  COALESCE(role_name,'') AS role_name, is_required
			 FROM ho_class_document ORDER BY id`)
	} else {
		err = r.db.SelectContext(ctx, &items,
			`SELECT id, ho_class_id, doc_class_id,
			  COALESCE(role_name,'') AS role_name, is_required
			 FROM ho_class_document WHERE ho_class_id = $1 ORDER BY id`, hoClassID)
	}
	return items, err
}

func (r *HoClassDocumentPostgres) Create(ctx context.Context, hcd *domain.HoClassDocument) (string, error) {
	var id string
	query := `INSERT INTO ho_class_document (ho_class_id, doc_class_id, role_name, is_required)
	          VALUES ($1, $2, NULLIF($3,''), $4)
	          RETURNING id`
	err := r.db.QueryRowContext(ctx, query, hcd.HoClassID, hcd.DocClassID, hcd.RoleName, hcd.IsRequired).Scan(&id)
	return id, err
}

func (r *HoClassDocumentPostgres) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM ho_class_document WHERE id = $1`, id)
	return err
}
