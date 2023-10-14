package handlers

import (
	"auth-crud/db"
	"auth-crud/handlers/auth"
	"auth-crud/models"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func handlePostQuery(query string, args ...interface{}) ([]models.Post, error) {
	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func handleSinglePostQuery(query string, args ...interface{}) (models.Post, error) {
	var post models.Post
	err := db.DB.QueryRow(query, args...).Scan(&post.ID, &post.Title, &post.Content, &post.UserID)
	return post, err
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	currentUserID := auth.GetCurrentUserIDFromContext(r)

	if currentUserID == 0 {
		http.Error(w, "Unauthorized", http.StatusForbidden)
	}
	// Use currentUserID to associate the post with the current user
	var newPost models.Post
	if err := json.NewDecoder(r.Body).Decode(&newPost); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newPost.UserID = currentUserID

	if _, err := db.DB.Exec("INSERT INTO posts (title, content, user_id) VALUES ($1, $2, $3)", newPost.Title, newPost.Content, newPost.UserID); err != nil {
		http.Error(w, "Error creating post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func ListPostHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := handlePostQuery("SELECT id, title, content, user_id FROM posts")
	if err != nil {
		http.Error(w, "Error querying posts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func GetPostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]

	post, err := handleSinglePostQuery("SELECT id, title, content, user_id FROM posts WHERE id = $1", postID)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func UpdatePostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]

	var updatedPost models.Post
	if err := json.NewDecoder(r.Body).Decode(&updatedPost); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, err := db.DB.Exec("UPDATE posts SET title = $1, content = $2 WHERE id = $3", updatedPost.Title, updatedPost.Content, postID); err != nil {
		http.Error(w, "Error updating post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID := vars["id"]

	if _, err := db.DB.Exec("DELETE FROM posts WHERE id = $1", postID); err != nil {
		http.Error(w, "Error deleting post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
