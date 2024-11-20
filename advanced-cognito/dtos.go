package main

type RegisterRequest struct {
	Email    string 
	Password string 
}

type VerificationCodeRequest struct {
	Email string 
	Code  string 
}

type LoginRequest struct {
	Email    string 
	Password string 
}

type UserCookie struct {
	Email        string
	AccessToken  string
	RefreshToken string
}
