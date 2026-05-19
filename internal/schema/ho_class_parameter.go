package schema

type CreateHoClassParameterInput struct {
	HoClassID   string  `json:"ho_class_id"  binding:"required"`
	ParameterID string  `json:"parameter_id" binding:"required"`
	OrderNum    int     `json:"order_num"`
	MinVal      float64 `json:"min_val"`
	MaxVal      float64 `json:"max_val"`
}

type UpdateHoClassParameterInput struct {
	OrderNum int     `json:"order_num"`
	MinVal   float64 `json:"min_val"`
	MaxVal   float64 `json:"max_val"`
}

type CopyHoClassParametersInput struct {
	FromClassID string `json:"from_class_id" binding:"required"`
	ToClassID   string `json:"to_class_id"   binding:"required"`
}
