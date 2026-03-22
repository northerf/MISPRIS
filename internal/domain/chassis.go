package domain

import "time"

type Chassis struct {
	ID            string `json:"chassis_id"      db:"chassis_id"`
	FrameID       string `json:"frame_id"        db:"frame_id"`
	SuspensionID  string `json:"suspension_id" db:"suspension_id"`
	BreakSystemID string `json:"break_system_id" db:"break_system_id"`
}

type Frame struct {
	ID        string    `json:"frame_id"        db:"frame_id"`
	Name      string    `json:"frame_name"      db:"frame_name"`
	Info      string    `json:"frame_info"      db:"frame_info"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt bool      `json:"deleted_at" db:"deleted_at"`
}

type Suspension struct {
	ID        string    `json:"suspension_id"        db:"suspension_id"`
	Name      string    `json:"suspension_name"      db:"suspension_name"`
	Info      string    `json:"suspension_info"      db:"suspension_info"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt bool      `json:"deleted_at" db:"deleted_at"`
}

type BreakSystem struct {
	ID        string    `json:"break_system_id"        db:"break_system_id"`
	Name      string    `json:"break_system_name"      db:"break_system_name"`
	Info      string    `json:"break_system_info"      db:"break_system_info"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt bool      `json:"deleted_at" db:"deleted_at"`
}
