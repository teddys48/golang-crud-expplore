package entity

type RoleDetail struct {
	ID     int64  `gorm:"column:id;primaryKey"`
	RoleID int64  `gorm:"role_id"`
	MenuID int64  `gorm:"menu_id"`
	Action string `gorm:"action"`
}
