package models

type User struct {
	UserInfo
	Password string `json:"password"`
}

type UserInfo struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	IsMember bool   `json:"is_chirpy_red"`
}

type UserLoginResponse struct {
	UserInfo
	Token   string `json:"token"`
	Refresh string `json:"refresh_token"`
}
