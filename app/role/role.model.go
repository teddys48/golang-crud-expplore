package role

import "time"

type Role struct {
	ID          int64            `json:"id"`
	Code        string           `json:"code"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Status      string           `json:"status"`
	CreatedBy   int64            `json:"created_by"`
	CreatedOn   time.Time        `json:"created_on"`
	UpdatedBy   int64            `json:"updated_by"`
	UpdatedOn   time.Time        `json:"updated_on"`
	MenuList    []RoleDetailData `json:"menu_list"`
}

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

type RoleDetailData struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Sort   int    `json:"sort"`
	Action string `json:"action"`
}

type RoleCreateRequest struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Menu        []RoleCreateRequestMenu `json:"menu"`
}

type RoleCreateRequestMenu struct {
	MenuID int64  `json:"menu_id"`
	Action string `json:"action"`
}
