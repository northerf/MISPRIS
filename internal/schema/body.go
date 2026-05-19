package schema

type CreateBodyInput struct {
	CarcassID string `json:"carcass_id"`
	DoorsID   string `json:"doors_id"`
	WingsID   string `json:"wings_id"`
}
type CreateCarcassInput struct {
	Name string `json:"name"`
	Info string `json:"info"`
}
type CreateDoorsInput struct {
	Name string `json:"name"`
	Info string `json:"info"`
}
type CreateWingsInput struct {
	Name string `json:"name"`
	Info string `json:"info"`
}
