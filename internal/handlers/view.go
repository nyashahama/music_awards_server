package handlers

import "net/http"

// ViewHandler handles nominee presentation
type ViewHandler interface {
	ListNominees(w http.ResponseWriter, r *http.Request)
	SearchNominees(w http.ResponseWriter, r *http.Request)
	GetPopularNominees(w http.ResponseWriter, r *http.Request)
	GetNomineeDetail(w http.ResponseWriter, r *http.Request)
	TrackNomineeView(w http.ResponseWriter, r *http.Request)
}

type viewHandler struct {
}