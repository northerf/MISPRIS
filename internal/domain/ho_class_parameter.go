package domain

// HoClassParameter — параметр, привязанный к типу ХО.
type HoClassParameter struct {
	ID          string  `json:"id"                  db:"id"`
	HoClassID   string  `json:"ho_class_id"         db:"ho_class_id"`
	ParameterID string  `json:"parameter_id"        db:"parameter_id"`
	OrderNum    int     `json:"order_num"           db:"order_num"`
	MinVal      float64 `json:"min_val"             db:"min_val"`
	MaxVal      float64 `json:"max_val"             db:"max_val"`
}

// HoClassParameterFull — результат SQL-функции get_ho_class_parameters.
// Содержит полные данные параметра вместе с ограничениями класса ХО.
type HoClassParameterFull struct {
	CpID          string  `json:"cp_id"          db:"cp_id"`
	ParamID       string  `json:"param_id"       db:"param_id"`
	Designation   string  `json:"designation"    db:"designation"`
	Name          string  `json:"name"           db:"name"`
	ParamType     string  `json:"param_type"     db:"param_type"`
	MeasuringUnit string  `json:"measuring_unit" db:"measuring_unit"`
	MinVal        float64 `json:"min_val"        db:"min_val"`
	MaxVal        float64 `json:"max_val"        db:"max_val"`
	OrderNum      int     `json:"order_num"      db:"order_num"`
}
