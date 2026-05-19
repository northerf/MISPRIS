package repository

import (
	"context"
	"fmt"
	"time"

	"MISPRIS/internal/domain"

	"github.com/jmoiron/sqlx"
)

type HoInstancePostgres struct {
	db *sqlx.DB
}

func NewHoInstancePostgres(db *sqlx.DB) *HoInstancePostgres {
	return &HoInstancePostgres{db: db}
}

// List возвращает экземпляры ХО; если hoClassID не пустой — фильтрует по типу.
func (r *HoInstancePostgres) List(ctx context.Context, hoClassID string) ([]*domain.HoInstance, error) {
	var items []*domain.HoInstance
	var err error
	if hoClassID == "" {
		err = r.db.SelectContext(ctx, &items,
			`SELECT ho_id, ho_class_id,
			  COALESCE(doc_number,'') AS doc_number,
			  COALESCE(doc_date::text,'') AS doc_date,
			  COALESCE(total_amount,0) AS total_amount,
			  COALESCE(status,'') AS status,
			  COALESCE(note,'') AS note,
			  created_at, updated_at
			 FROM ho_instance ORDER BY ho_id`)
	} else {
		err = r.db.SelectContext(ctx, &items,
			`SELECT ho_id, ho_class_id,
			  COALESCE(doc_number,'') AS doc_number,
			  COALESCE(doc_date::text,'') AS doc_date,
			  COALESCE(total_amount,0) AS total_amount,
			  COALESCE(status,'') AS status,
			  COALESCE(note,'') AS note,
			  created_at, updated_at
			 FROM ho_instance WHERE ho_class_id = $1 ORDER BY ho_id`, hoClassID)
	}
	return items, err
}

func (r *HoInstancePostgres) GetByID(ctx context.Context, id string) (*domain.HoInstance, error) {
	var item domain.HoInstance
	query := `SELECT ho_id, ho_class_id,
	           COALESCE(doc_number,'') AS doc_number,
	           COALESCE(doc_date::text,'') AS doc_date,
	           COALESCE(total_amount,0) AS total_amount,
	           COALESCE(status,'') AS status,
	           COALESCE(note,'') AS note,
	           created_at, updated_at
	          FROM ho_instance WHERE ho_id = $1`
	err := r.db.GetContext(ctx, &item, query, id)
	return &item, err
}

// Create вызывает SQL-функцию ins_ho, возвращающую INTEGER (новый ho_id).
func (r *HoInstancePostgres) Create(ctx context.Context, h *domain.HoInstance) (string, error) {
	var id string
	query := `SELECT ins_ho($1, NULLIF($2,''), NULLIF($3,'')::date, NULLIF($4,0), NULLIF($5,''))`
	err := r.db.QueryRowContext(ctx, query,
		h.HoClassID, h.DocNumber, h.DocDate, h.TotalAmount, h.Note,
	).Scan(&id)
	return id, err
}

func (r *HoInstancePostgres) Update(ctx context.Context, h *domain.HoInstance) error {
	query := `UPDATE ho_instance
	          SET status=COALESCE(NULLIF($1,''), status),
	              doc_number=COALESCE(NULLIF($2,''), doc_number),
	              doc_date=COALESCE(NULLIF($3,'')::date, doc_date),
	              total_amount=COALESCE(NULLIF($4,0), total_amount),
	              note=COALESCE(NULLIF($5,''), note),
	              updated_at=$6
	          WHERE ho_id = $7`
	_, err := r.db.ExecContext(ctx, query,
		h.Status, h.DocNumber, h.DocDate, h.TotalAmount, h.Note, time.Now(), h.ID)
	return err
}

func (r *HoInstancePostgres) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM ho_instance WHERE ho_id = $1`, id)
	return err
}

// FindByClass — прямой запрос вместо SQL-функции.
// FIX: функция возвращала INTEGER ho_id и TEXT actors, несовместимые с полями struct.
func (r *HoInstancePostgres) FindByClass(ctx context.Context, hoClassID string) ([]*domain.HoInstanceFull, error) {
	var items []*domain.HoInstanceFull
	query := `
		SELECT
			hi.ho_id::text                                              AS ho_id,
			COALESCE(hi.doc_number, '')                              AS doc_number,
			COALESCE(hi.doc_date::text, '')                          AS doc_date,
			COALESCE(hi.total_amount, 0)                              AS total_amount,
			COALESCE(hi.status, '')                                  AS status,
			COUNT(DISTINCT hp.id)                                      AS positions,
			COALESCE(
				STRING_AGG(DISTINCT s.name || '   ('  || hr.name || ')'  , ', '),
				''
			)                                                          AS actors
		FROM ho_instance hi
		LEFT JOIN ho_position hp  ON hp.ho_id  = hi.ho_id
		LEFT JOIN ho_actor    ha  ON ha.ho_id  = hi.ho_id
		LEFT JOIN shd         s   ON s.shd_id  = ha.shd_id
		LEFT JOIN ho_role     hr  ON hr.ho_role_id = ha.ho_role_id
		WHERE hi.ho_class_id = $1
		GROUP BY hi.ho_id
		ORDER BY hi.doc_date DESC, hi.ho_id DESC`
	err := r.db.SelectContext(ctx, &items, query, hoClassID)
	if err != nil {
		return nil, fmt.Errorf("FindByClass: %w", err)
	}
	return items, nil
}
