package main

import (
	"Ketab/core"
	"net/http"
)

func main() {

	handler := core.NewHandler()
	http.HandleFunc("/register", handler.Register)
	http.HandleFunc("/login", handler.Login)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
