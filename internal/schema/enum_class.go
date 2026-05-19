package schema

type CreateEnumClassInput struct {
	Name          string `json:"name"           binding:"required"`
	ComponentType string `json:"component_type"`
}

type UpdateEnumClassInput struct {
	Name          string `json:"name"`
	ComponentType string `json:"component_type"`
}

type ValidateEnumValueInput struct {
	EnumClassID string `json:"enum_class_id" binding:"required"`
	Value       string `json:"value"         binding:"required"`
}
