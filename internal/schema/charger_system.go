package schema

type CreateChargerSystemInput struct {
	ChargerID   string `json:"charger_id"`
	ConnectorID string `json:"connector_id"`
}
type CreateChargerInput struct {
	Name string `json:"name"`
	Info string `json:"info"`
}
type CreateConnectorInput struct {
	Name string `json:"name"`
	Info string `json:"info"`
}
