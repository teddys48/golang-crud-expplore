package auth

import "github.com/golang-jwt/jwt/v5"

type LoginRequest struct {
	EmailOrNIP string `json:"email_or_nip" validate:"required"`
	Password   string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	Menu         *[]Menu `json:"menu"`
	User         any     `json:"user"`
}

type ClaimsToken struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type LoginUsers struct {
	ID       int64  `json:"id"`
	Code     string `json:"code"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	RoleID   int64  `json:"role_id"`
}

type UsersData struct {
	// Username     string `json:"username"`
	Email        string  `json:"email"`
	NIP          *string `json:"nip"`
	NRK          string  `json:"nrk"`
	PerusahaanID int     `json:"perusahaan_id"`
	Instansi     string  `json:"instansi"`
	Jabatan      string  `json:"jabatan"`
	Phone        string  `json:"phone"`
	Fullname     string  `json:"fullname"`
	Role         string  `json:"role"`
}

type Menu struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Sort   int    `json:"sort"`
	Action string `json:"action"`
}
