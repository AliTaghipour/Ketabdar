package main

import (
	"Ketab/core"
	"net/http"
)

func main() {

	handler := core.NewHandler()
	http.HandleFunc("/register", handler.Register)
	http.HandleFunc("/login", handler.Login)
	http.HandleFunc("/book", handler.AddBook)
	http.HandleFunc("/my_books", handler.GetBooks)
	http.HandleFunc("/update_book", handler.UpdateBook)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
