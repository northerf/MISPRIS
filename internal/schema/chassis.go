package schema

type CreateChassisInput struct {
	FrameID       string `json:"frame_id"`
	SuspensionID  string `json:"suspension_id"`
	BreakSystemID string `json:"break_system_id"`
}
type CreateFrameInput struct {
	Name string `json:"name"`
	Info string `json:"info"`
}
type CreateSuspensionInput struct {
	Name string `json:"name"`
	Info string `json:"info"`
}
type CreateBreakSystemInput struct {
	Name string `json:"name"`
	Info string `json:"info"`
}
