package schema

type CreateEnumPositionInput struct {
	EnumClassID string `json:"enum_class_id" binding:"required"`
	Value       string `json:"value"         binding:"required"`
	OrderNum    string `json:"order_num"`
}

type UpdateEnumPositionInput struct {
	Value    string `json:"value"`
	OrderNum string `json:"order_num"`
}

type ReorderEnumPositionInput struct {
	NewOrderNum string `json:"new_order_num" binding:"required"`
}
