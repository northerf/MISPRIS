package domain

import "time"

type ChargerSystem struct {
	ID          string    `json:"charger_system_id" db:"charger_system_id"`
	ChargerID   string    `json:"charger_id"        db:"charger_id"`
	ConnectorID string    `json:"connector_id" db:"connector_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt   bool      `json:"deleted_at" db:"deleted_at"`
}

type Charger struct {
	ID        string    `json:"charger_id"        db:"charger_id"`
	Name      string    `json:"charger_name"      db:"charger_name"`
	Info      string    `json:"charger_info"      db:"charger_info"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt bool      `json:"deleted_at" db:"deleted_at"`
}

type Connector struct {
	ID        string    `json:"connector_id"        db:"connector_id"`
	Name      string    `json:"connector_name"      db:"connector_name"`
	Info      string    `json:"connector_info"      db:"connector_info"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt bool      `json:"deleted_at" db:"deleted_at"`
}
