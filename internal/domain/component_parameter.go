package domain

// ComponentParameter — привязка параметра к типу компонента с ограничениями.
//
// ComponentType — тип компонента: "emobile", "battery", "engine", "chassis", ...
// MinVal / MaxVal — ограничения для числовых параметров в контексте данного типа.
//
// Пример: для «emobile» параметр «Запас хода» ограничен диапазоном 50–800 км.
type ComponentParameter struct {
	ID            string  `json:"component_parameter_id" db:"component_parameter_id"`
	ComponentType string  `json:"component_type"         db:"component_type"`
	ParameterID   string  `json:"parameter_id"           db:"parameter_id"`
	OrderNum      int     `json:"order_num"              db:"order_num"`
	MinVal        float64 `json:"min_val"                db:"min_val"`
	MaxVal        float64 `json:"max_val"                db:"max_val"`
}

// ComponentParameterFull — результат SQL-функции get_component_parameters.
// Содержит полные данные о параметре вместе с ограничениями класса.
type ComponentParameterFull struct {
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
