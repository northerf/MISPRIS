package schema

type CreateComponentParameterInput struct {
	ComponentType string `json:"component_type" binding:"required"`
	ParameterID   string `json:"parameter_id"   binding:"required"`
	OrderNum      string `json:"order_num"`
	MinVal        string `json:"min_val"`
	MaxVal        string `json:"max_val"`
}

type UpdateComponentParameterInput struct {
	OrderNum string `json:"order_num"`
	MinVal   string `json:"min_val"`
	MaxVal   string `json:"max_val"`
}

type CopyComponentParametersInput struct {
	FromType string `json:"from_type" binding:"required"`
	ToType   string `json:"to_type"   binding:"required"`
}
