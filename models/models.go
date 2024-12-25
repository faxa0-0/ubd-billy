package models

import "time"

type Role string

const (
	Customer Role = "customer"
	Admin    Role = "admin"
)

type User struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Login  string `json:"login"`
	Pass   string `json:"pass,omitempty"`
	Plan   string `json:"plan"`
	Role   Role   `json:"role"`
	Active bool   `json:"active"`
}

type Usage struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id,omitempty"`
	Youtube    uint64    `json:"youtube"` //kb
	Netflix    uint64    `json:"netflix"`
	Spotify    uint64    `json:"spotify"`
	Basic      uint64    `json:"basic"`
	VerifiedAt time.Time `json:"verified_at"`
}
