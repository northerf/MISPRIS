package domain

import "time"

type Battery struct {
	ID              string    `json:"battery_id"    db:"battery_id"`
	BatteryName     string    `json:"battery_name" db:"battery_name"`
	BatteryType     string    `json:"battery_type"   db:"battery_type"`
	BatteryCapacity string    `json:"battery_capacity" db:"battery_capacity"`
	BatteryInfo     string    `json:"battery_info"   db:"battery_info"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt       bool      `json:"deleted_at" db:"deleted_at"`
}
