package models

// Status represents the status of attendance
type Status struct {
	ID   int    `json:"id"`
	Statusname string `json:"status_name"`
	Code string `json:"code"`
}
