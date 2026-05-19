package repository

import (
	"context"

	"MISPRIS/internal/domain"

	"github.com/jmoiron/sqlx"
)

type HoDocumentPostgres struct {
	db *sqlx.DB
}

func NewHoDocumentPostgres(db *sqlx.DB) *HoDocumentPostgres {
	return &HoDocumentPostgres{db: db}
}

// ListByHo возвращает документы ХО; если hoID не пустой — фильтрует по экземпляру ХО.
func (r *HoDocumentPostgres) ListByHo(ctx context.Context, hoID string) ([]*domain.HoDocument, error) {
	var items []*domain.HoDocument
	var err error
	if hoID == "" {
		err = r.db.SelectContext(ctx, &items,
			`SELECT id, ho_id, doc_class_id,
			  COALESCE(doc_number,'') AS doc_number,
			  COALESCE(doc_date::text,'') AS doc_date,
			  COALESCE(note,'') AS note
			 FROM ho_document ORDER BY id`)
	} else {
		err = r.db.SelectContext(ctx, &items,
			`SELECT id, ho_id, doc_class_id,
			  COALESCE(doc_number,'') AS doc_number,
			  COALESCE(doc_date::text,'') AS doc_date,
			  COALESCE(note,'') AS note
			 FROM ho_document WHERE ho_id = $1 ORDER BY id`, hoID)
	}
	return items, err
}

func (r *HoDocumentPostgres) Create(ctx context.Context, d *domain.HoDocument) (string, error) {
	var id string
	query := `INSERT INTO ho_document (ho_id, doc_class_id, doc_number, doc_date, note)
	          VALUES ($1, $2, NULLIF($3,''), NULLIF($4,'')::date, NULLIF($5,''))
	          RETURNING id`
	err := r.db.QueryRowContext(ctx, query,
		d.HoID, d.DocClassID, d.DocNumber, d.DocDate, d.Note,
	).Scan(&id)
	return id, err
}

func (r *HoDocumentPostgres) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM ho_document WHERE id = $1`, id)
	return err
}
