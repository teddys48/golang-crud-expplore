package jobs

import (
	"encoding/json"
	"time"
)

type Jobs struct {
	ID         int64           `json:"id"`
	Code       string          `json:"code"`
	Name       string          `json:"name"`
	ProjectID  int64           `json:"project_id"`
	Data       json.RawMessage `json:"data"`
	Notes      string          `json:"notes"`
	ApprovedBy *int64          `json:"approved_by"`
	ApprovedOn *time.Time      `json:"approved_on"`
	Status     string          `json:"status"`
	CreatedBy  int64           `json:"created_by"`
	CreatedOn  time.Time       `json:"created_on"`
	UpdatedBy  *int64          `json:"updated_by"`
	UpdatedOn  *time.Time      `json:"updated_on"`
}
