package models

import "time"

type Transaction struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	CategoryID  int       `json:"category_id"`
	Category    string    `json:"category,omitempty"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
