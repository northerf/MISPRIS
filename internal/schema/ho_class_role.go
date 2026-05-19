package schema

type CreateHoClassRoleInput struct {
	HoClassID  string `json:"ho_class_id" binding:"required"`
	HoRoleID   string `json:"ho_role_id"  binding:"required"`
	IsRequired bool   `json:"is_required"`
}
