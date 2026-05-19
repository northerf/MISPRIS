package schema

type CreateHoRoleInput struct {
	Name        string `json:"name"        binding:"required"`
	Description string `json:"description"`
}

type UpdateHoRoleInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
