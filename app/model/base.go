package model

import (
	"time"
)

// Base defines the basic fields for all models
type Base struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
