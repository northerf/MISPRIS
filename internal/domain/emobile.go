package domain

import "time"

type Emobile struct {
	ID              string    `json:"emobile_id"        db:"emobile_id"`
	Name            string    `json:"emobile_name"       db:"emobile_name"`
	PowerPointID    string    `json:"power_point_id"     db:"power_point_id"`
	BatteryID       string    `json:"battery_id"         db:"battery_id"`
	ChargerSystemID string    `json:"charger_system_id"  db:"charger_system_id"`
	ChassisID       string    `json:"chassis_id"         db:"chassis_id"`
	BodyID          string    `json:"body_id"            db:"body_id"`
	ElectronicsID   string    `json:"electronics_id"     db:"electronics_id"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt       bool      `json:"deleted_at" db:"deleted_at"`
}
