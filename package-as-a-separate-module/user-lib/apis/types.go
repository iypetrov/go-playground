package apis

import (
	go_yaml "github.com/goccy/go-yaml"
)

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func NewUser(id int64, username, email string) User {
	return User{
		ID:       id,
		Username: username,
		Email:    email,
	}
}

func (u *User) ToYaml() string {
	bytes, err := go_yaml.Marshal(*u)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
