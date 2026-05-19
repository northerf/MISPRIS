package domain

import "time"

// HoClass — тип (класс) хозяйственной операции (классификатор).
type HoClass struct {
	ID         string    `json:"ho_class_id"  db:"ho_class_id"`
	Name       string    `json:"name"         db:"name"`
	Designation string   `json:"designation"  db:"designation"`
	ParentID   string    `json:"parent_id"    db:"parent_id"`
	IsTerminal bool      `json:"is_terminal"  db:"is_terminal"`
	CreatedAt  time.Time `json:"created_at"   db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"   db:"updated_at"`
}
