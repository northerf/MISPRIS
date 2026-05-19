package schema

type CreateShdInput struct {
	Name        string `json:"name"        binding:"required"`
	ShdType     string `json:"shd_type"`
	INN         string `json:"inn"`
	Description string `json:"description"`
}

type UpdateShdInput struct {
	Name        string `json:"name"`
	ShdType     string `json:"shd_type"`
	INN         string `json:"inn"`
	Description string `json:"description"`
}
