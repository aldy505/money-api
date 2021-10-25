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
	StatusRequested Status = iota
	StatusRejected
	StatusPending
	StatusSuccess
	StatusCancelled
	StatusFailed
)

type Intermmediate struct {
	Sender    int
	Recipient int
	Amount    int64
	Message   string
	Status    Status
}
