package model

import (
	"time"
)

type Order struct {
	ID         int       `gorm:"type:integer;primary_key" json:"id,omitempty"`
	Price      int       `gorm:"type:integer;not null" json:"price,omitempty"`
	UserId     int       `gorm:"type:integer;not null" json:"user_id,omitempty"`
	Status     string    `json:"status"`
	Reason     string    `json:"reason"`
	CreatedAt  time.Time `gorm:"not null" json:"created_at,omitempty"`
	ModifiedAt time.Time `gorm:"not null" json:"modified_at,omitempty"`
}

type NewOrder struct {
	UserId          int    `json:"user_id,omitempty"`
	Price           int    `json:"price,omitempty"`
	DeliveryAddress string `json:"delivery_address"`
	DeliveryDate    string `json:"delivery_date"`
}
