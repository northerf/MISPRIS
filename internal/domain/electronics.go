package domain

import "time"

type Electronics struct {
	ID           string    `json:"electronics_id"  db:"electronics_id"`
	ControllerID string    `json:"controller_id"   db:"controller_id"`
	SensorID     string    `json:"sensor_id"       db:"sensor_id"`
	WiringID     string    `json:"wiring_id"       db:"wiring_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt    bool      `json:"deleted_at" db:"deleted_at"`
}

type Controller struct {
	ID        string    `json:"controller_id"   db:"controller_id"`
	Name      string    `json:"controller_name" db:"controller_name"`
	Info      string    `json:"controller_info" db:"controller_info"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt bool      `json:"deleted_at" db:"deleted_at"`
}

type Sensor struct {
	ID        string    `json:"sensor_id"   db:"sensor_id"`
	Name      string    `json:"sensor_name" db:"sensor_name"`
	Info      string    `json:"sensor_info" db:"sensor_info"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt bool      `json:"deleted_at" db:"deleted_at"`
}

type Wiring struct {
	ID        string    `json:"wiring_id"   db:"wiring_id"`
	Name      string    `json:"wiring_name" db:"wiring_name"`
	Info      string    `json:"wiring_info" db:"wiring_info"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt bool      `json:"deleted_at" db:"deleted_at"`
}
