package auth

import "github.com/golang-jwt/jwt/v5"

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Menu         *[]string `json:"menu"`
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
}
