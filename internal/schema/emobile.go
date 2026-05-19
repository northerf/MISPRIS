package schema

type CreateEmobileInput struct {
	Name            string `json:"name" binding:"required"`
	PowerPointID    string `json:"power_point_id" binding:"required"`
	BatteryID       string `json:"battery_id" binding:"required"`
	ChargerSystemID string `json:"charger_system_id" binding:"required"`
	ChassisID       string `json:"chassis_id" binding:"required"`
	BodyID          string `json:"body_id" binding:"required"`
	ElectronicsID   string `json:"electronics_id" binding:"required"`
}

type UpdateEmobileInput struct {
	Name            string `json:"name"`
	PowerPointID    string `json:"power_point_id"`
	BatteryID       string `json:"battery_id"`
	ChargerSystemID string `json:"charger_system_id"`
	ChassisID       string `json:"chassis_id"`
	BodyID          string `json:"body_id"`
	ElectronicsID   string `json:"electronics_id"`
}
