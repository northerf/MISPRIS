package domain

import "time"

// HoInstance — экземпляр хозяйственной операции (конкретный документ/сделка).
type HoInstance struct {
	ID          string    `json:"ho_id"         db:"ho_id"`
	HoClassID   string    `json:"ho_class_id"   db:"ho_class_id"`
	DocNumber   string    `json:"doc_number"    db:"doc_number"`
	DocDate     string    `json:"doc_date"      db:"doc_date"`
	TotalAmount float64   `json:"total_amount"  db:"total_amount"`
	Status      string    `json:"status"        db:"status"`
	Note        string    `json:"note"          db:"note"`
	CreatedAt   time.Time `json:"created_at"    db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"    db:"updated_at"`
}

// HoInstanceFull — результат запроса FindByClass (сводная информация об операциях).
// Actors — строка вида «Завод (Отправитель), ООО Ромашка (Получатель)».
type HoInstanceFull struct {
	HoID        string  `json:"ho_id"         db:"ho_id"`
	DocNumber   string  `json:"doc_number"    db:"doc_number"`
	DocDate     string  `json:"doc_date"      db:"doc_date"`
	TotalAmount float64 `json:"total_amount"  db:"total_amount"`
	Status      string  `json:"status"        db:"status"`
	Positions   int64   `json:"positions"     db:"positions"`
	Actors      string  `json:"actors"        db:"actors"` // было int — баг!
}
