package auth

import "time"

type User struct {
	ID        int       `json:"id,omitempty" db:"id"`
	Name      string    `json:"name" db:"name"`
	Password  string    `json:"password,omitempty" db:"password"`
	Email     string    `json:"email,omitempty" db:"email"`
	Address   string    `json:"address,omitempty" db:"address"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
