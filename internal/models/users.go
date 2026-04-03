package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // hide this
	Role      string    `json:"role"`
	IsActive  bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
}


type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}


type LoginInput struct{
	Email string `json: "email"`
	Password string `json: "password`
}

