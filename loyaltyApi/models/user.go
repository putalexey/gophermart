package models

var UserIdentityKey = "UUID"

type User struct {
	UUID     string
	Login    string
	Password string
}
