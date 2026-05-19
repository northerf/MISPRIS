package domain

// HoRole — роль участника хозяйственной операции (покупатель, продавец, перевозчик и т.д.)
type HoRole struct {
	ID          string `json:"ho_role_id"   db:"ho_role_id"`
	Name        string `json:"name"         db:"name"`
	Description string `json:"description"  db:"description"`
}
