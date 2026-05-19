package schema

type CreatePowerPointInput struct {
	EngineID   string `json:"engine_id"`
	InverterID string `json:"inverter_id"`
	GearboxID  string `json:"gearbox_id"`
}
type CreateEngineInput struct {
	Name       string `json:"name"`
	EngineType string `json:"engine_type"`
	Info       string `json:"info"`
}
type CreateInverterInput struct {
	Name string `json:"name"`
	Info string `json:"info"`
}
type CreateGearboxInput struct {
	Name string `json:"name"`
	Info string `json:"info"`
}
