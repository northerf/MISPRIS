package domain

import "time"

// EnumPosition — одно значение внутри справочника перечисления.
// Например, для «Стандарт зарядки»: «CCS2», «CHAdeMO», «Type 2».
// order_num управляет порядком отображения в списке.
type EnumPosition struct {
	ID          string    `json:"enum_position_id" db:"enum_position_id"`
	EnumClassID string    `json:"enum_class_id"    db:"enum_class_id"`
	Value       string    `json:"value"            db:"value"`
	OrderNum    int       `json:"order_num"        db:"order_num"`
	CreatedAt   time.Time `json:"created_at"       db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"       db:"updated_at"`
}
