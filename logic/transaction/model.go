package transaction

import (
	"money-api/logic/account"
	"time"
)

type Transaction struct {
	ID        int             `json:"id" db:"id"`
	Sender    account.Account `json:"from"`
	Recipient account.Account `json:"to"`
	Amount    int64           `json:"amount" db:"amount"`
	Message   string          `json:"message,omitempty" db:"message"`
	Status    Status          `json:"status" db:"status"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
}

type Status int

const (
	StatusRequested Status = iota + 1
	StatusRejected
	StatusPending
	StatusSuccess
	StatusCancelled
	StatusFailed
)

type Intermediate struct {
	Sender    int    `json:"sender_id"`
	Recipient int    `json:"recipient_id"`
	Amount    int64  `json:"amount"`
	Message   string `json:"message"`
	Status    Status `json:"status"`
}
