package models

type User struct {
	Id int64
	Name string
	Points int64
}

type Users []*User