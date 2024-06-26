package model

type RequestRegister struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type RequestLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
