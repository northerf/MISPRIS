package domain

// HoDocument — документ, прикреплённый к экземпляру хозяйственной операции.
type HoDocument struct {
	ID         string `json:"id"           db:"id"`
	HoID       string `json:"ho_id"        db:"ho_id"`
	DocClassID string `json:"doc_class_id" db:"doc_class_id"`
	DocNumber  string `json:"doc_number"   db:"doc_number"`
	DocDate    string `json:"doc_date"     db:"doc_date"`
	Note       string `json:"note"         db:"note"`
}
