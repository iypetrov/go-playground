package apis

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
