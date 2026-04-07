package main

import (
	"fmt"

	"github.com/iypetrov/user-lib/apis"
)

func main() {
	users := []apis.User{
		apis.NewUser(1, "John Smith", "john.smith@gmail.com"),
		apis.NewUser(2, "John Adams", "john.admams@gmail.com"),
		apis.NewUser(3, "John Black", "john.black@gmail.com"),
	}
	for _, u := range users {
		fmt.Println(u.ToYaml())
	}
}
