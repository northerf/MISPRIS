package schema

type CreateEmobileParameterValueInput struct {
	EmobileID            string `json:"emobile_id"              binding:"required"`
	ComponentParameterID string `json:"component_parameter_id"  binding:"required"`
	ValReal              string `json:"val_real"`
	ValInt               string `json:"val_int"`
	ValStr               string `json:"val_str"`
	EnumValID            string `json:"enum_val_id"`
}

type UpdateEmobileParameterValueInput struct {
	ValReal   string `json:"val_real"`
	ValInt    string `json:"val_int"`
	ValStr    string `json:"val_str"`
	EnumValID string `json:"enum_val_id"`
}
