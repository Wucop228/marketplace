package models

type User struct {
	ID           int64
	Username     string
	Role         string
	PasswordHash string
}
