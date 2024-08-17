package corporation

import "time"

type Corporation struct {
	ID            int64      `json:"id"`
	Code          string     `json:"code"`
	Name          *string    `json:"name"`
	Address       *string    `json:"address"`
	Npwp          *string    `json:"npwp"`
	DirectorName  *string    `json:"director_name"`
	Email         *string    `json:"email"`
	Fax           *string    `json:"fax"`
	NotarisNumber *string    `json:"notaris_number"`
	NotarisDate   *string    `json:"notaris_date"`
	Status        string     `json:"status"`
	CreatedBy     int64      `json:"created_by"`
	CreatedOn     time.Time  `json:"created_on"`
	UpdatedBy     *int64     `json:"updated_by"`
	UpdatedOn     *time.Time `json:"updated_on"`
}
