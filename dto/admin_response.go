package dto

type AdminResponse struct {
	Id       string `json:"id_admin"`
	Nik      string `json:"nik"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
