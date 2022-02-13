package models

var UserIdentityKey = "UUID"

type User struct {
	UUID     string `json:"uuid"`
	Login    string `json:"login"`
	Password string `json:"-"`
}
