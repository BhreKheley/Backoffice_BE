package models

// Role represents the role of a user in the system
type Role struct {
	ID   int    `json:"id"`
	Rolename string `json:"role_name"`
	Code string `json:"code"`
}
