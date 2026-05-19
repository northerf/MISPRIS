package domain

import "time"

// EnumClass — класс (справочник) перечисления.
// Хранит допустимые значения для атрибута компонента.
// Например: «Стандарт зарядки», «Класс защиты IP», «Тип привода».
type EnumClass struct {
	ID            string    `json:"enum_class_id"  db:"enum_class_id"`
	Name          string    `json:"name"           db:"name"`
	ComponentType string    `json:"component_type" db:"component_type"`
	CreatedAt     time.Time `json:"created_at"     db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"     db:"updated_at"`
}
