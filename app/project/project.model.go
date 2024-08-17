package project

import "time"

type Project struct {
	ID                            int64      `json:"id"`
	Code                          string     `json:"code"`
	ActivityName                  string     `json:"activity_name"`
	PackageName                   string     `json:"package_name"`
	AccountCode                   string     `json:"account_code"`
	AccountDescription            string     `json:"account_description"`
	CorporationId                 int64      `json:"corporation_id"`
	Executor                      string     `json:"executor"`
	Npwp                          string     `json:"npwp"`
	Director                      string     `json:"director"`
	Position                      string     `json:"position"`
	Email                         string     `json:"email"`
	Fax                           string     `json:"fax"`
	NotarisNumber                 string     `json:"notaris_number"`
	NotarisDate                   string     `json:"notaris_date"`
	NotarisName                   string     `json:"notaris_name"`
	RupID                         string     `json:"rup_id"`
	CatalogID                     string     `json:"catalog_id"`
	SppbjNumber                   string     `json:"sppbj_number"`
	SppbjDate                     string     `json:"sppbj_date"`
	ContractNumber                string     `json:"contract_number"`
	ContractDate                  string     `json:"contract_date"`
	ContractDateNumber            string     `json:"contract_date_number"`
	ContractDateLetter            string     `json:"contract_date_letter"`
	ContractMonth                 string     `json:"contract_month"`
	ContractDay                   string     `json:"contract_day"`
	ContractAddNumber             string     `json:"contract_add_number"`
	ContractAddDate               string     `json:"contract_add_date"`
	ContractDateAddNumber         string     `json:"contract_date_add_number"`
	ContractDateAddLetter         string     `json:"contract_date_add_letter"`
	ContractAddMonth              string     `json:"contract_add_month"`
	ContractAddDay                string     `json:"contract_add_day"`
	OrderLetterNumber             string     `json:"order_letter_number"`
	SpDate                        string     `json:"sp_date"`
	SpmkNumber                    string     `json:"spmk_number"`
	SpmkDate                      string     `json:"spmk_date"`
	ImplementationScheduleDay     string     `json:"implementation_schedule_day"`
	Start                         string     `json:"start"`
	End                           string     `json:"end"`
	Kelurahan                     string     `json:"kelurahan"`
	Kecamatan                     string     `json:"kecamatan"`
	SppbjValueLetter              string     `json:"sppbj_value_letter"`
	SppbjValue                    int64      `json:"sppbj_value"`
	NilaiEPurchasingOk            int64      `json:"nilai_e_purchasing_ok"`
	NilaiEPurchasing              int64      `json:"nilai_e_purchasing"`
	NilaiPembulatan               int64      `json:"nilai_pembulatan"`
	ContractValue                 int64      `json:"contract_value"`
	ContractValueLetter           string     `json:"contract_value_letter"`
	BusinessClassification        string     `json:"business_classification"`
	Volume                        int64      `json:"volume"`
	Unit                          string     `json:"unit"`
	AddendumType                  string     `json:"addendum_type"`
	TargetPercentage              int64      `json:"target_percentage"`
	ProgressPlan                  int64      `json:"progress_plan"`
	TotalCumulativeProgressPlan   int64      `json:"total_cumulative_progress_plan"`
	ActualProgress                int64      `json:"actual_progress"`
	TotalCumulativeActualProgress int64      `json:"total_cumulative_actual_progress"`
	Deviation                     int64      `json:"deviation"`
	Status                        string     `json:"status"`
	CreatedBy                     int64      `json:"created_by"`
	CreatedOn                     time.Time  `json:"created_on"`
	UpdatedBy                     *int64     `json:"updated_by"`
	UpdatedOn                     *time.Time `json:"updated_on"`
}
