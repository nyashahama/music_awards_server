package handlers

import (
	"net/http"

	"github.com/nyashahama/music-awards/internal/services"
)

// UserHandler handles user-related HTTP requests
type UserHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	GetProfile(w http.ResponseWriter, r *http.Request)
	UpdateProfile(w http.ResponseWriter, r *http.Request)
	DeleteAccount(w http.ResponseWriter, r *http.Request)
	PromoteUser(w http.ResponseWriter, r *http.Request)
}


type userHandler struct {
	userService services.UserService
}
