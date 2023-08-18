package model

import (
	"time"
)

type Billing struct {
	ID         int       `gorm:"type:integer;primary_key" json:"id,omitempty"`
	UserId     int       `gorm:"not null" json:"user_id"`
	Balance    int       `gorm:"type:integer" json:"balance"`
	CreatedAt  time.Time `gorm:"not null" json:"created_at,omitempty"`
	ModifiedAt time.Time `gorm:"not null" json:"modified_at,omitempty"`
}
