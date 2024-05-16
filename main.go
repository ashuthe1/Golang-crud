package main

import (
	"crud/controller"
	"crud/helper"
	"log"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/api/books", booksHandler)
	http.HandleFunc("/api/books/create", createBookHandler)
	http.HandleFunc("/api/book/", bookHandler)

	config := helper.GetConfiguration()

	log.Println("Server started on", config.Port)

	err := http.ListenAndServe(config.Port, nil)
	if err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}

func createBookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	controller.CreateBook(w, r)
}
func booksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		controller.GetAllBooks(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func bookHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/book/")
	if id == "" || strings.Contains(id, "/") {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case "GET":
		controller.GetBook(w, r, id)
	case "PUT":
		controller.UpdateBook(w, r, id)
	case "DELETE":
		controller.DeleteBook(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
