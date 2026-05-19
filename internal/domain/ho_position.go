package domain

// HoPosition — позиция (строка) хозяйственной операции, содержащая ссылку на ТС.
type HoPosition struct {
	ID          string  `json:"id"            db:"id"`
	HoID        string  `json:"ho_id"         db:"ho_id"`
	EmobileID   string  `json:"emobile_id"    db:"emobile_id"`
	Quantity    int     `json:"quantity"      db:"quantity"`
	UnitPrice   float64 `json:"unit_price"    db:"unit_price"`
	TotalPrice  float64 `json:"total_price"   db:"total_price"`
	Note        string  `json:"note"          db:"note"`
	PositionNum int     `json:"position_num"  db:"position_num"`
}

// HoPositionFull — позиция ХО с присоединённым названием ТС.
type HoPositionFull struct {
	ID          string  `json:"id"            db:"id"`
	HoID        string  `json:"ho_id"         db:"ho_id"`
	EmobileID   string  `json:"emobile_id"    db:"emobile_id"`
	EmobileName string  `json:"emobile_name"  db:"emobile_name"`
	Quantity    int     `json:"quantity"      db:"quantity"`
	UnitPrice   float64 `json:"unit_price"    db:"unit_price"`
	TotalPrice  float64 `json:"total_price"   db:"total_price"`
	Note        string  `json:"note"          db:"note"`
	PositionNum int     `json:"position_num"  db:"position_num"`
}
