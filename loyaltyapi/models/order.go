package models

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	UUID       string    `json:"uuid"`
	UserUUID   string    `json:"user_uuid"`
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    int       `json:"accrual"`
	UploadedAt time.Time `json:"uploaded_at"`
}

func NewOrder() *Order {
	return &Order{
		UUID:       uuid.NewString(),
		Status:     OrderStatusRegistered,
		UploadedAt: time.Now(),
	}
}

const (
	OrderStatusRegistered = "REGISTERED"
	OrderStatusProcessed  = "PROCESSED"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusInvalid    = "INVALID"
)
