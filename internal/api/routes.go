package api

import (
	"github.com/Ammar022/evv-logger-backend/internal/repo"
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router, store repo.Store) {
	handler := NewHandler(store)

	// Schedule routes
	router.HandleFunc("/schedules", handler.handleGetSchedules).Methods("GET")
	router.HandleFunc("/schedules/{id}", handler.handleGetScheduleByID).Methods("GET")
	router.HandleFunc("/schedules/{id}/start", handler.handleStartVisit).Methods("POST")
	router.HandleFunc("/schedules/{id}/end", handler.handleEndVisit).Methods("POST")
	
	// Task routes
	router.HandleFunc("/tasks/{taskId}/update", handler.handleUpdateTask).Methods("POST")
	
	// Stats route for dashboard
	router.HandleFunc("/stats", handler.handleGetStats).Methods("GET")
}
