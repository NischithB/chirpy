package models

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserInfo struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}
