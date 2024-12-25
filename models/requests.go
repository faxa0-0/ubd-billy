package models

type LoginRequest struct {
	Login string `json:"login"`
	Pass  string `json:"pass"`
}
