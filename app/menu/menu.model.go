package menu

import "time"

type Menu struct {
	ID           int64     `json:"id"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	MenuParentID int64     `json:"menu_parent_id"`
	Icon         string    `json:"icon"`
	PathURL      string    `json:"path_url"`
	Sort         string    `json:"sort"`
	HiddenData   bool      `json:"hidden_data"`
	Description  string    `json:"description"`
	Status       string    `json:"status"`
	CreatedBy    int64     `json:"created_by"`
	CreatedOn    time.Time `json:"created_on"`
	UpdatedBy    int64     `json:"updated_by"`
	UpdatedOn    time.Time `json:"updated_on"`
}

type MenuCreateRequest struct {
	Name         string `json:"name"`
	MenuParentID int64  `json:"menu_parent_id"`
	Icon         string `json:"icon"`
	PathURL      string `json:"path_url"`
	Sort         int    `json:"sort"`
	HiddenData   bool   `json:"hidden_data"`
	Description  string `json:"description"`
}

type MenuUpdateRequest struct {
	Name         string `json:"name"`
	MenuParentID int64  `json:"menu_parent_id"`
	Icon         string `json:"icon"`
	PathURL      string `json:"path_url"`
	Sort         int    `json:"sort"`
	HiddenData   bool   `json:"hidden_data"`
	Description  string `json:"description"`
}
