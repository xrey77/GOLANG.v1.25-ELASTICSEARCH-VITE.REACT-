package models

import (
	"time"
)

type User struct {
	ID          string    `json:"id"`
	Lastname    string    `json:"lastname"`
	Firstname   string    `json:"firstname"`
	Email       string    `json:"email"`
	Mobile      string    `json:"mobile"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	Roles       string    `json:"roles"`
	Isactivated bool      `json:"isactivated"`
	Isblocked   bool      `json:"isblocked"`
	Userpicture string    `json:"userpicture"`
	Mailtoken   float64   `json:"mailtoken"`
	Secret      *string   `json:"secret"`
	Qrcodeurl   *string   `json:"qrcodeurl"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
