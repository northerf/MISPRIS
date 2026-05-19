package domain

// HoParameterValue — значение параметра для конкретного экземпляра ХО.
type HoParameterValue struct {
	ID                 string  `json:"id"                    db:"id"`
	HoID               string  `json:"ho_id"                 db:"ho_id"`
	HoClassParameterID string  `json:"ho_class_parameter_id" db:"ho_class_parameter_id"`
	ValReal            float64 `json:"val_real"              db:"val_real"`
	ValInt             int     `json:"val_int"               db:"val_int"`
	ValStr             string  `json:"val_str"               db:"val_str"`
	ValDate            string  `json:"val_date"              db:"val_date"`
	EnumValID          string  `json:"enum_val_id"           db:"enum_val_id"`
}

// HoParameterValueFull — значение параметра с именем и типом из таблицы parameter.
// Используется при ListByHo чтобы фронтенд мог показать человекочитаемое название.
type HoParameterValueFull struct {
	ID                 string  `json:"id"                    db:"id"`
	HoID               string  `json:"ho_id"                 db:"ho_id"`
	HoClassParameterID string  `json:"ho_class_parameter_id" db:"ho_class_parameter_id"`
	ParamName          string  `json:"param_name"            db:"param_name"`
	ParamType          string  `json:"param_type"            db:"param_type"`
	MeasuringUnit      string  `json:"measuring_unit"        db:"measuring_unit"`
	ValReal            float64 `json:"val_real"              db:"val_real"`
	ValInt             int     `json:"val_int"               db:"val_int"`
	ValStr             string  `json:"val_str"               db:"val_str"`
	ValDate            string  `json:"val_date"              db:"val_date"`
	EnumValID          string  `json:"enum_val_id"           db:"enum_val_id"`
}
