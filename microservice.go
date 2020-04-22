package main

import (
	"net/http"
	"os"

	"cloud-native-go/api"
)

func main() {
	http.HandleFunc("/api/books", api.BooksHandeFunc)
	http.HandleFunc("/api/books/", api.BookHandeFunc)

	if err := http.ListenAndServe(port(), nil); err != nil {
		panic(err)
	}
}

func port() string {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	return ":" + port
}
