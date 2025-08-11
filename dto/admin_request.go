package dto

type AdminRequest struct {
	Nik      string `json:"nik"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Nik      string `json:"nik"`
	Password string `json:"password"`
}