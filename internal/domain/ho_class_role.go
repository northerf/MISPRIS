package domain

// HoClassRole — связь типа ХО с обязательной или необязательной ролью.
type HoClassRole struct {
	ID         string `json:"id"           db:"id"`
	HoClassID  string `json:"ho_class_id"  db:"ho_class_id"`
	HoRoleID   string `json:"ho_role_id"   db:"ho_role_id"`
	IsRequired bool   `json:"is_required"  db:"is_required"`
}
