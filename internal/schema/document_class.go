package schema

type CreateDocumentClassInput struct {
	Name        string `json:"name"        binding:"required"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

type UpdateDocumentClassInput struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
}
