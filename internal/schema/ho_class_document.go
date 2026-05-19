package schema

type CreateHoClassDocumentInput struct {
	HoClassID  string `json:"ho_class_id"  binding:"required"`
	DocClassID string `json:"doc_class_id" binding:"required"`
	RoleName   string `json:"role_name"`
	IsRequired bool   `json:"is_required"`
}
