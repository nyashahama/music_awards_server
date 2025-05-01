package handlers

import "net/http"

// NotificationHandler handles notification endpoints
type NotificationHandler interface {
	SendVotingReminders(w http.ResponseWriter, r *http.Request)
	GetNotifications(w http.ResponseWriter, r *http.Request)
	MarkNotificationRead(w http.ResponseWriter, r *http.Request)
}

type notificationHandler struct {
	// Dependencies
}