package models

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	UUID       string    `json:"uuid"`
	UserUUID   string    `json:"user_uuid" db:"user_uuid"`
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    int       `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at" db:"uploaded_at"`
}

func NewOrder() *Order {
	return &Order{
		UUID:       uuid.NewString(),
		Status:     OrderStatusNew,
		UploadedAt: time.Now(),
	}
}

const (
	OrderStatusNew        = "NEW"
	OrderStatusProcessed  = "PROCESSED"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusInvalid    = "INVALID"
)
