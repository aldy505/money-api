package account

import "time"

type Account struct {
	ID        int       `json:"id,omitempty" db:"id"`
	Balance   int64     `json:"balance,omitempty" db:"balance"`
	Tag       string    `json:"tag,omitempty" db:"tag"`
	Friends   []Friend  `json:"friend,omitempty"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Friend struct {
	Friend int `json:"id" db:"friend"`
}
