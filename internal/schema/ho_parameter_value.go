package schema

type CreateHoParameterValueInput struct {
	HoID               string  `json:"ho_id"                 binding:"required"`
	HoClassParameterID string  `json:"ho_class_parameter_id" binding:"required"`
	ValReal            float64 `json:"val_real"`
	ValInt             int     `json:"val_int"`
	ValStr             string  `json:"val_str"`
	ValDate            string  `json:"val_date"`
	EnumValID          string  `json:"enum_val_id"`
}

type UpdateHoParameterValueInput struct {
	ValReal   float64 `json:"val_real"`
	ValInt    int     `json:"val_int"`
	ValStr    string  `json:"val_str"`
	ValDate   string  `json:"val_date"`
	EnumValID string  `json:"enum_val_id"`
}
