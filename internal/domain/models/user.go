package models

type User struct {
	ID       int
	Name     string
	Username string
	PassHash []byte
}
