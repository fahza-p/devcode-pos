package models

type Pagination struct {
	TotalData int64 `json:"total"`
	Limit     int   `json:"limit"`
	Skip      int   `json:"skip"`
}
