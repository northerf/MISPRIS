package schema

type CreateElectronicsInput struct {
	ControllerID string `json:"controller_id"`
	SensorID     string `json:"sensor_id"`
	WiringID     string `json:"wiring_id"`
}
type CreateControllerInput struct {
	Name string `json:"name"`
	Info string `json:"info"`
}
type CreateSensorInput struct {
	Name string `json:"name"`
	Info string `json:"info"`
}
type CreateWiringInput struct {
	Name string `json:"name"`
	Info string `json:"info"`
}
