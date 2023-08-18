package model

import (
	"time"
)

type BillingAccount struct {
	ID         int       `gorm:"type:integer;primary_key" json:"id,omitempty"`
	UserId     int       `gorm:"not null" json:"user_id"`
	Balance    int       `gorm:"type:integer" json:"balance"`
	CreatedAt  time.Time `gorm:"not null" json:"created_at,omitempty"`
	ModifiedAt time.Time `gorm:"not null" json:"modified_at,omitempty"`
}

type NewAccount struct {
	ID         int       `json:"id,omitempty"`
	UserId     int       `json:"user_id"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	ModifiedAt time.Time `json:"modified_at,omitempty"`
}

type Deposit struct {
	UserId int `json:"user_id"`
	Amount int `json:"amount"`
}

type Withdrawal struct {
	UserId int `json:"user_id"`
	Amount int `json:"amount"`
}
