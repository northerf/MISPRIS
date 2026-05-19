package domain

// EmobileParameterValue — значение параметра для конкретного автомобиля.
//
// В зависимости от типа параметра заполняется одно из val_*:
//   - ValReal  → для "real" (запас хода 320.5)
//   - ValInt   → для "int"  (масса 1850)
//   - ValStr   → для "str"  (страна «Германия»)
//   - EnumValID → для "enum" (ссылка на EnumPosition)
type EmobileParameterValue struct {
	ID                    string  `json:"value_id"                db:"value_id"`
	EmobileID             string  `json:"emobile_id"              db:"emobile_id"`
	ComponentParameterID  string  `json:"component_parameter_id"  db:"component_parameter_id"`
	ValReal               float64 `json:"val_real"                db:"val_real"`
	ValInt                int     `json:"val_int"                 db:"val_int"`
	ValStr                string  `json:"val_str"                 db:"val_str"`
	EnumValID             string  `json:"enum_val_id"             db:"enum_val_id"`
}
