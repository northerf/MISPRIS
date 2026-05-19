package schema

type CreateHoPositionInput struct {
	HoID      string  `json:"ho_id"      binding:"required"`
	EmobileID string  `json:"emobile_id" binding:"required"`
	Quantity  int     `json:"quantity"   binding:"required"`
	UnitPrice float64 `json:"unit_price"`
	Note      string  `json:"note"`
}

type UpdateHoPositionInput struct {
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
	Note      string  `json:"note"`
}
