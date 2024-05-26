package models

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type UserSignInUp struct {
	Username     string
	PasswordHash string
}

type UserCreatedInfo struct {
	ID       int64
	Username string
}
