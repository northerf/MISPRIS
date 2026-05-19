package schema

type CreateHoDocumentInput struct {
	HoID       string `json:"ho_id"        binding:"required"`
	DocClassID string `json:"doc_class_id" binding:"required"`
	DocNumber  string `json:"doc_number"`
	DocDate    string `json:"doc_date"`
	Note       string `json:"note"`
}
