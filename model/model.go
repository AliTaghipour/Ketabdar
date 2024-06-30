package model

type BookStatus int

const (
	BookStatusAlreadyRead BookStatus = iota
	BookStatusReading
	BookStatusPending
)

type User struct {
	Id       int32
	Username string
	Password string
	Email    string
}

type Book struct {
	Id     int32
	Title  string
	Author string
	Pages  int32
}

type UserBook struct {
	Book
	Status    BookStatus
	ReadPages int32
}
