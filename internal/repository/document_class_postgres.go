package repository

import (
	"context"

	"MISPRIS/internal/domain"

	"github.com/jmoiron/sqlx"
)

type DocumentClassPostgres struct {
	db *sqlx.DB
}

func NewDocumentClassPostgres(db *sqlx.DB) *DocumentClassPostgres {
	return &DocumentClassPostgres{db: db}
}

func (r *DocumentClassPostgres) List(ctx context.Context) ([]*domain.DocumentClass, error) {
	var items []*domain.DocumentClass
	query := `SELECT doc_class_id, name,
	           COALESCE(code,'') AS code,
	           COALESCE(description,'') AS description,
	           created_at
	          FROM document_class ORDER BY doc_class_id`
	err := r.db.SelectContext(ctx, &items, query)
	return items, err
}

func (r *DocumentClassPostgres) GetByID(ctx context.Context, id string) (*domain.DocumentClass, error) {
	var item domain.DocumentClass
	query := `SELECT doc_class_id, name,
	           COALESCE(code,'') AS code,
	           COALESCE(description,'') AS description,
	           created_at
	          FROM document_class WHERE doc_class_id = $1`
	err := r.db.GetContext(ctx, &item, query, id)
	return &item, err
}

func (r *DocumentClassPostgres) Create(ctx context.Context, dc *domain.DocumentClass) (string, error) {
	var id string
	query := `INSERT INTO document_class (name, code, description, created_at)
	          VALUES ($1, NULLIF($2,''), NULLIF($3,''), NOW())
	          RETURNING doc_class_id`
	err := r.db.QueryRowContext(ctx, query, dc.Name, dc.Code, dc.Description).Scan(&id)
	return id, err
}

func (r *DocumentClassPostgres) Update(ctx context.Context, dc *domain.DocumentClass) error {
	query := `UPDATE document_class SET name=$1, code=NULLIF($2,''), description=NULLIF($3,'')
	          WHERE doc_class_id = $4`
	_, err := r.db.ExecContext(ctx, query, dc.Name, dc.Code, dc.Description, dc.ID)
	return err
}

func (r *DocumentClassPostgres) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM document_class WHERE doc_class_id = $1`, id)
	return err
}
