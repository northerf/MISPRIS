package domain

import "time"

// SHD — субъект хозяйственной деятельности (организация, ИП и т.п.)
type SHD struct {
	ID          string    `json:"shd_id"       db:"shd_id"`
	Name        string    `json:"name"         db:"name"`
	ShdType     string    `json:"shd_type"     db:"shd_type"`
	INN         string    `json:"inn"          db:"inn"`
	Description string    `json:"description"  db:"description"`
	CreatedAt   time.Time `json:"created_at"   db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"   db:"updated_at"`
}
