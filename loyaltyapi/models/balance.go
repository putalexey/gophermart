package models

import "time"

type Balance struct {
	UserUUID  string  `json:"-" db:"user_uuid"`
	Current   float64 `json:"current" db:"current"`
	Withdrawn float64 `json:"withdrawn" db:"withdrawn"`
}

type Withdrawal struct {
	UUID        string    `json:"-" db:"uuid"`
	UserUUID    string    `json:"-" db:"user_uuid"`
	Order       string    `json:"order" db:"order"`
	Sum         float64   `json:"sum" db:"sum"`
	ProcessedAt time.Time `json:"processed_at" db:"processed_at"`
}
