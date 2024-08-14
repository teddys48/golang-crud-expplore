package entity

import "time"

type Role struct {
	ID          int64     `gorm:"column:id;primaryKey"`
	Code        string    `gorm:"code"`
	Name        string    `gorm:"name"`
	Description string    `gorm:"description"`
	Status      string    `gorm:"status"`
	CreatedBy   int64     `gorm:"created_by"`
	CreatedOn   time.Time `gorm:"created_on"`
	UpdatedBy   int64     `gorm:"updated_by"`
	UpdatedOn   time.Time `gorm:"updated_on"`
}
