package schema

type CreateHoClassInput struct {
	Name        string `json:"name"        binding:"required"`
	Designation string `json:"designation"`
	ParentID    string `json:"parent_id"`
	IsTerminal  bool   `json:"is_terminal"`
}

type UpdateHoClassInput struct {
	Name        string `json:"name"`
	Designation string `json:"designation"`
	ParentID    string `json:"parent_id"`
	IsTerminal  bool   `json:"is_terminal"`
}
