package main

import (
	"fmt"

	"github.com/iypetrov/user-lib/apis"
)

func main() {
	users := []apis.User{
		apis.NewUser(1, "John Smith", "john.smith@gmail.com"),
	}
	for _, u := range users {
		fmt.Println(u.ID)
	}
}
