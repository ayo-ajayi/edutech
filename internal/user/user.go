package user

import (
	"time"
)

type User struct {
	Email      string    `json:"email" bson:"email"`
	Password   string    `json:"-" bson:"password"`
	Firstname  string    `json:"firstname" bson:"firstname"`
	Lastname   string    `json:"lastname" bson:"lastname"`
	IsVerified bool      `json:"is_verified" bson:"is_verified"`
	Role       Role      `json:"role" bson:"role"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" bson:"updated_at"`
}
type Role string

const (
	Admin   Role = "admin"
	Tutor   Role = "tutor"
	Student Role = "student"
)
