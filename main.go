package main

import (
	"Ketab/core"
	"net/http"
)

func main() {

	handler := core.NewHandler()
	http.HandleFunc("/register", handler.Register)
	http.HandleFunc("/login", handler.Login)
	http.HandleFunc("/add_book", handler.AddBook)
	http.HandleFunc("/books", handler.GetBooks)
	http.HandleFunc("/lib/add_book", handler.AddBookToLibrary)
	http.HandleFunc("/lib/books", handler.GetLibraryBooks)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
