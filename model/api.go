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

type RequestAddBook struct {
	Book Book `json:"book"`
}

type RequestAddBookLibrary struct {
	UserBook UserBook `json:"user_book"`
}
