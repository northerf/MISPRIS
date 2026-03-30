package domain

import "time"

type User struct {
	ID        string    `json:"user_id" db:"user_id"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"-" db:"password"`
	Role      string    `json:"role" db:"role"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
