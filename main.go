package main

import (
	"auth-crud/db"
	"auth-crud/handlers"

	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	db.InitDB() // Initialize your database connection

	r := mux.NewRouter()

	// Create a subrouter for the "posts" resource
	v1 := r.PathPrefix("/api/v1").Subrouter()

	v1.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	v1.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	v1.HandleFunc("/logout", handlers.LogoutHandler).Methods("POST")

	v1.HandleFunc("/posts", handlers.CreatePostHandler).Methods("POST")
	v1.HandleFunc("/posts", handlers.ListPostHandler).Methods("GET")
	v1.HandleFunc("/post/{id}", handlers.GetPostHandler).Methods("GET")
	v1.HandleFunc("/post/{id}", handlers.UpdatePostHandler).Methods("PUT")
	v1.HandleFunc("/post/{id}", handlers.DeletePostHandler).Methods("DELETE")

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
