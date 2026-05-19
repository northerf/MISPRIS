package schema

type CreateHoActorInput struct {
	HoID     string `json:"ho_id"      binding:"required"`
	HoRoleID string `json:"ho_role_id" binding:"required"`
	ShdID    string `json:"shd_id"     binding:"required"`
}
