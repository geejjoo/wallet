package task

import "time"

type Wallet struct {
	ID      string  `json:"id"`
	Balance float64 `json:"balance"`
}

type SendRequest struct {
	ToID   string  `json:"to"`
	Amount float64 `json:"amount"`
}

type History struct {
	FromID string    `json:"from_id" db:"from_id"`
	ToID   string    `json:"to_id" db:"to_id"`
	Amount float64   `json:"amount" db:"amount"`
	Time   time.Time `json:"time" db:"time"`
}

type UpdateWallet struct {
	FromID string  `json:"from_id" db:"from_id"`
	ToID   string  `json:"to_id" db:"to_id"`
	Amount float64 `json:"amount" db:"amount"`
}
