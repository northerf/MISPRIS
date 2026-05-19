package domain

// HoActor — участник (актор) хозяйственной операции в заданной роли.
type HoActor struct {
	ID       string `json:"id"          db:"id"`
	HoID     string `json:"ho_id"       db:"ho_id"`
	HoRoleID string `json:"ho_role_id"  db:"ho_role_id"`
	ShdID    string `json:"shd_id"      db:"shd_id"`
}
