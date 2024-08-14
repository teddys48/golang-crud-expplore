package entity

import "time"

type Menu struct {
	ID           int64     `gorm:"column:id;primaryKey"`
	Code         string    `gorm:"code"`
	Name         string    `gorm:"name"`
	MenuParentID int64     `gorm:"menu_parent_id"`
	Icon         string    `gorm:"icon"`
	PathURL      string    `gorm:"path_url"`
	Sort         int       `gorm:"sort"`
	HiddenData   bool      `gorm:"hidden_data"`
	Description  string    `gorm:"description"`
	Status       string    `gorm:"status"`
	CreatedBy    int64     `gorm:"created_by"`
	CreatedOn    time.Time `gorm:"created_on"`
	UpdatedBy    int64     `gorm:"updated_by"`
	UpdatedOn    time.Time `gorm:"updated_on"`
}
