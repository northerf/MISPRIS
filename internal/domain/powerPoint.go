package domain

import "time"

type PowerPoint struct {
	ID         string    `json:"power_point_id" db:"power_point_id"`
	EngineID   string    `json:"engine_id"        db:"engine_id"`
	InverterID string    `json:"inverter_id" db:"inverter_id"`
	GearboxID  string    `json:"gearbox_id" db:"gearbox_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt  bool      `json:"deleted_at" db:"deleted_at"`
}

type Engine struct {
	ID         string    `json:"engine_id" db:"engine_id"`
	EngineName string    `json:"engine_name" db:"engine_name"`
	EngineType string    `json:"engine_type" db:"engine_type"`
	EngineInfo string    `json:"engine_info" db:"engine_info"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt  bool      `json:"deleted_at" db:"deleted_at"`
}

type Inverter struct {
	ID           string    `json:"inverter_id" db:"inverter_id"`
	InverterName string    `json:"inverter_name" db:"inverter_name"`
	InverterInfo string    `json:"inverter_info" db:"inverter_info"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt    bool      `json:"deleted_at" db:"deleted_at"`
}

type Gearbox struct {
	ID          string    `json:"gearbox_id" db:"gearbox_id"`
	GearboxName string    `json:"gearbox_name" db:"gearbox_name"`
	GearboxInfo string    `json:"gearbox_info" db:"gearbox_info"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt   bool      `json:"deleted_at" db:"deleted_at"`
}
