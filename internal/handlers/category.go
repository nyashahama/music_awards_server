package handlers

import "net/http"

// CategoryHandler handles category operations
type CategoryHandler interface {
	CreateCategory(w http.ResponseWriter, r *http.Request)
	UpdateCategory(w http.ResponseWriter, r *http.Request)
	DeleteCategory(w http.ResponseWriter, r *http.Request)
	GetCategory(w http.ResponseWriter, r *http.Request)
	ListCategories(w http.ResponseWriter, r *http.Request)
}

type categoryHandler struct {
	// Dependencies
}