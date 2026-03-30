package schema

type CreateChargerSystemInput struct {
	ChargerID   string `json:"charger_id" binding:"required"`
	ConnectorID string `json:"connector_id" binding:"required"`
}

type CreateChargerInput struct {
	Name string `json:"name" binding:"required"`
	Info string `json:"info"`
}

type CreateConnectorInput struct {
	Name string `json:"name" binding:"required"`
	Info string `json:"info"`
}
