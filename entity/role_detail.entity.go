package entity

import "time"

type RoleDetail struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	RoleID    int64     `gorm:"role_id"`
	MenuID    int64     `gorm:"menu_id"`
	Action    string    `gorm:"action"`
	Status    string    `gorm:"status"`
	CreatedBy int64     `gorm:"created_by"`
	CreatedOn time.Time `gorm:"created_on"`
	UpdatedBy int64     `gorm:"updated_by"`
	UpdatedOn time.Time `gorm:"updated_on"`
}
