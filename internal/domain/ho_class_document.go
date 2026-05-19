package domain

// HoClassDocument — требование к документу для данного типа ХО.
type HoClassDocument struct {
	ID         string `json:"id"           db:"id"`
	HoClassID  string `json:"ho_class_id"  db:"ho_class_id"`
	DocClassID string `json:"doc_class_id" db:"doc_class_id"`
	RoleName   string `json:"role_name"    db:"role_name"`
	IsRequired bool   `json:"is_required"  db:"is_required"`
}
