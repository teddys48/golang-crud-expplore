package users

import "time"

type UserCreateRequest struct {
	Username     string  `json:"username"`
	Email        string  `json:"email" validate:"required"`
	NIP          *string `json:"nip"`
	NRK          *string `json:"nrk"`
	PerusahaanID *int    `json:"perusahaan_id"`
	Instansi     *string `json:"instansi"`
	Jabatan      *string `json:"jabatan"`
	Phone        *string `json:"phone"`
	Fullname     *string `json:"fullname"`
	Password     string  `json:"password" validate:"required"`
	RoleID       int     `json:"role_id" validate:"required"`
}

type UserUpdateRequest struct {
	Username     string  `json:"username"`
	Email        string  `json:"email" validate:"required"`
	NIP          *string `json:"nip"`
	NRK          *string `json:"nrk"`
	PerusahaanID *int    `json:"perusahaan_id"`
	Instansi     *string `json:"instansi"`
	Jabatan      *string `json:"jabatan"`
	Phone        *string `json:"phone"`
	Fullname     *string `json:"fullname"`
	RoleID       int     `json:"role_id" validate:"required"`
}

type Users struct {
	ID            int64     `json:"id"`
	Code          string    `json:"code"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	NIP           string    `json:"nip"`
	CorporationID int64     `json:"corporation_id"`
	Instance      string    `json:"instance"`
	Position      string    `json:"position"`
	NRK           string    `json:"nrk"`
	Phone         string    `json:"phone"`
	Password      string    `json:"password"`
	Fullname      string    `json:"full_name"`
	ProfileImage  string    `json:"profile_image"`
	RoleID        int64     `json:"role_id"`
	Status        string    `json:"status"`
	CreatedBy     int64     `json:"created_by"`
	CreatedOn     time.Time `json:"created_on"`
	UpdatedBy     int64     `json:"updated_by"`
	UpdatedOn     time.Time `json:"updated_on"`
}
