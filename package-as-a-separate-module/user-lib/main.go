package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/iypetrov/user-lib/apihelpers"
	"github.com/iypetrov/user-lib/apis"
)

func usersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	users := []apis.User{
		apis.NewUser(1, "John Smith", "john.smith@gmail.com"),
		apis.NewUser(2, "John Adams", "john.admams@gmail.com"),
		apis.NewUser(3, "John Black", "john.black@gmail.com"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
	for _, u := range users {
		fmt.Println(apihelpers.ToYaml(u))
	}
}

func main() {
	http.HandleFunc("/users", usersHandler)

	fmt.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
