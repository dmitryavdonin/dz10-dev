package model

import (
	"time"
)

type User struct {
	ID         int       `gorm:"type:integer;primary_key" json:"id,omitempty"`
	Username   string    `gorm:"not null" json:"username"`
	Balance    int       `gorm:"type:integer" json:"balance"`
	CreatedAt  time.Time `gorm:"not null" json:"created_at,omitempty"`
	ModifiedAt time.Time `gorm:"not null" json:"modified_at,omitempty"`
}
