package domain

import "time"

// DocumentClass — класс (тип) документа, сопровождающего хозяйственную операцию.
type DocumentClass struct {
	ID          string    `json:"doc_class_id" db:"doc_class_id"`
	Name        string    `json:"name"         db:"name"`
	Code        string    `json:"code"         db:"code"`
	Description string    `json:"description"  db:"description"`
	CreatedAt   time.Time `json:"created_at"   db:"created_at"`
}
