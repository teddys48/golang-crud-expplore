package entity

import "time"

type Jobs struct {
	ID         int64       `gorm:"column:id;primaryKey"`
	Code       string      `gorm:"code"`
	Name       string      `gorm:"name"`
	ProjectID  int64       `gorm:"project_id"`
	Data       interface{} `gorm:"data"`
	Notes      string      `gorm:"notes"`
	ApprovedBy int64       `gorm:"approved_by"`
	ApprovedOn time.Time   `gorm:"approved_on"`
	Status     string      `gorm:"status"`
	CreatedBy  int64       `gorm:"created_by"`
	CreatedOn  time.Time   `gorm:"created_on"`
	UpdatedBy  int64       `gorm:"updated_by"`
	UpdatedOn  time.Time   `gorm:"updated_on"`
}
