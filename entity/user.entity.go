package entity

import "time"

type Users struct {
	ID            int64     `gorm:"column:id;primaryKey"`
	Code          string    `gorm:"code"`
	Username      string    `gorm:"username"`
	Email         string    `gorm:"email"`
	Nip           *string   `gorm:"nip" json:"nip"`
	CorporationID *int64    `gorm:"corporation_id"`
	Instance      *string   `gorm:"instance"`
	Position      *string   `gorm:"position"`
	NRK           *string   `gorm:"nrk"`
	Phone         *string   `gorm:"phone"`
	Password      string    `gorm:"password"`
	FullName      string    `gorm:"full_name"`
	ProfileImage  string    `gorm:"profile_image"`
	RoleID        int64     `gorm:"role_id"`
	Status        string    `gorm:"status"`
	CreatedBy     int64     `gorm:"created_by"`
	CreatedOn     time.Time `gorm:"created_on"`
	UpdatedBy     int64     `gorm:"updated_by"`
	UpdatedOn     time.Time `gorm:"updated_on"`
}
