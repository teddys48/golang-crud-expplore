package entity

import "time"

type Corporation struct {
	ID            int64     `gorm:"column:id;primaryKey"`
	Code          string    `gorm:"code"`
	Name          string    `gorm:"name"`
	Address       string    `gorm:"address"`
	Npwp          string    `gorm:"npwp"`
	DirectorName  string    `gorm:"director_name"`
	Email         string    `gorm:"email"`
	Fax           string    `gorm:"fax"`
	NotarisNumber string    `gorm:"notaris_number"`
	NotarisDate   time.Time `gorm:"notaris_date"`
	Status        string    `gorm:"status"`
	CreatedBy     int64     `gorm:"created_by"`
	CreatedOn     time.Time `gorm:"created_on"`
	UpdatedBy     int64     `gorm:"updated_by"`
	UpdatedOn     time.Time `gorm:"updated_on"`
}
