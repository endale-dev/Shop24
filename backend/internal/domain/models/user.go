package models

import "time"

type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	GoogleSub    string    `json:"-"`
	ResetHash    string    `json:"-"`
	ResetExpires time.Time `json:"-"`
	HasPassword  bool      `json:"hasPassword"`
	GoogleLinked bool      `json:"googleLinked"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
