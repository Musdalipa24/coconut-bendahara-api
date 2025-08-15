package dto

type AdminRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateAdminRequest struct {
	OldPassword string `json:"old_password"`
	Password    string `json:"password"`
}
