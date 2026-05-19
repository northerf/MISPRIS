package domain

import "time"

// Parameter — описание одной характеристики компонента/автомобиля.
//
// Типы параметров (ParamType):
//   - "real" — вещественное число (запас хода 320.5 км)
//   - "int"  — целое число (масса 1850 кг)
//   - "str"  — строка (страна сборки «Германия»)
//   - "enum" — значение из перечисления (стандарт зарядки «CCS2»)
//
// Для типа "enum" заполняется EnumClassID — ссылка на допустимые значения.
type Parameter struct {
	ID            string    `json:"parameter_id"   db:"parameter_id"`
	Designation   string    `json:"designation"    db:"designation"`
	Name          string    `json:"name"           db:"name"`
	ParamType     string    `json:"param_type"     db:"param_type"`
	MeasuringUnit string    `json:"measuring_unit" db:"measuring_unit"`
	EnumClassID   string    `json:"enum_class_id"  db:"enum_class_id"`
	CreatedAt     time.Time `json:"created_at"     db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"     db:"updated_at"`
}
