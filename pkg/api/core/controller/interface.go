package controller

import "time"

type Controller struct {
	Token1 string `json:"user_token"`
	Token2 string `json:"tmp_token"` //Hash(Token2 + Token3)
}

type Chat struct {
	ID        uint      `json:"id"`
	Err       string    `json:"error"`
	CreatedAt time.Time `json:"created_at"`
	UserID    uint      `json:"user_id"`
	GroupID   uint      `json:"group_id"`
	Admin     bool      `json:"admin"`
	Message   string    `json:"message"`
}