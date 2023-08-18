package model

import (
	"time"
)

type BillingTransaction struct {
	ID         int       `gorm:"type:integer;primary_key" json:"id,omitempty"`
	UserId     int       `json:"user_id"`
	OrderId    int       `json:"order_id"`
	Operation  string    `json:"operation"`
	Amount     int       `json:"amount"`
	Status     string    `json:"status"`
	Reason     string    `json:"reason"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}
