package models

import "time"

type Job struct {
	UUID      string    `json:"uuid"`
	OrderUUID string    `json:"order_uuid"`
	ProceedAt time.Time `json:"proceed_at"`
	Tries     int       `json:"tries"`
}
