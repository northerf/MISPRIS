package schema

type CreateParameterInput struct {
	Designation   string `json:"designation"    binding:"required"`
	Name          string `json:"name"           binding:"required"`
	ParamType     string `json:"param_type"     binding:"required"`
	MeasuringUnit string `json:"measuring_unit"`
	EnumClassID   string `json:"enum_class_id"`
}

type UpdateParameterInput struct {
	Designation   string `json:"designation"`
	Name          string `json:"name"`
	ParamType     string `json:"param_type"`
	MeasuringUnit string `json:"measuring_unit"`
	EnumClassID   string `json:"enum_class_id"`
}
