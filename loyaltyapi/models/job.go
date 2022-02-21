package models

import "time"

type Job struct {
	UUID      string    `json:"uuid" db:"uuid"`
	OrderUUID string    `json:"order_uuid" db:"order_uuid"`
	ProceedAt time.Time `json:"proceed_at" db:"proceed_at"`
	Tries     int       `json:"tries" db:"tries"`
}
