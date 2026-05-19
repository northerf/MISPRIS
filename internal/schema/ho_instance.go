package schema

type CreateHoInstanceInput struct {
	HoClassID   string  `json:"ho_class_id"  binding:"required"`
	DocNumber   string  `json:"doc_number"`
	DocDate     string  `json:"doc_date"`
	TotalAmount float64 `json:"total_amount"`
	Note        string  `json:"note"`
}

type UpdateHoInstanceInput struct {
	Status      string  `json:"status"`
	DocNumber   string  `json:"doc_number"`
	DocDate     string  `json:"doc_date"`
	TotalAmount float64 `json:"total_amount"`
	Note        string  `json:"note"`
}
