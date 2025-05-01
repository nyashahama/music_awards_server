package handlers

import "net/http"

// ResultsHandler handles results and reporting
type ResultsHandler interface {
	GetCategoryResults(w http.ResponseWriter, r *http.Request)
	GenerateReport(w http.ResponseWriter, r *http.Request)
	ExportResults(w http.ResponseWriter, r *http.Request)
	GetHistoricalResults(w http.ResponseWriter, r *http.Request)
	GetRealTimeTallies(w http.ResponseWriter, r *http.Request)
}

type resultsHandler struct {
	// Dependencies
}