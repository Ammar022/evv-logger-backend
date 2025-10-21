package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Ammar022/evv-logger-backend/internal/models"
	"github.com/Ammar022/evv-logger-backend/internal/repo"
	"github.com/gorilla/mux"
)

type Handler struct {
	store repo.Store
}

func NewHandler(store repo.Store) *Handler {
	return &Handler{store: store}
}

// @Summary Get all schedules
// @Description Retrieves a list of all caregiver schedules. Use the optional 'date' query parameter with the value 'today' to filter for schedules on the current day.
// @Tags schedules
// @Accept  json
// @Produce  json
// @Param date query string false "Filter by date" Enums(today) example(today)
// @Success 200 {array} models.Schedule "A list of schedules"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /schedules [get]
func (h *Handler) handleGetSchedules(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET /schedules - Request received")

	dateFilter := r.URL.Query().Get("date")

	var schedules []models.Schedule
	var err error

	if dateFilter == "today" {
		schedules, err = h.store.GetTodaySchedules()
		log.Printf("Fetching today's schedules")
	} else {
		schedules, err = h.store.GetSchedules()
		log.Printf("Fetching all schedules")
	}

	if err != nil {
		log.Printf("Error fetching schedules: %v", err)
		http.Error(w, "Failed to fetch schedules", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(schedules); err != nil {
		log.Printf("Error encoding schedules response: %v", err)
	}

	log.Printf("Successfully returned %d schedules", len(schedules))
}

// @Summary Get a single schedule by ID
// @Description Retrieves full details for a single schedule, including its assigned care activities and tasks.
// @Tags schedules
// @Accept  json
// @Produce  json
// @Param id path int true "Schedule ID" example(1)
// @Success 200 {object} models.Schedule "Full details of the schedule"
// @Failure 400 {object} map[string]string "Invalid Schedule ID format"
// @Failure 404 {object} map[string]string "Schedule not found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /schedules/{id} [get]
func (h *Handler) handleGetScheduleByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	log.Printf("GET /schedules/%s - Request received", idStr)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid schedule ID: %s", idStr)
		http.Error(w, "Invalid schedule ID", http.StatusBadRequest)
		return
	}

	schedule, err := h.store.GetScheduleByID(id)
	if err != nil {
		log.Printf("Error fetching schedule %d: %v", id, err)
		http.Error(w, "Schedule not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(schedule); err != nil {
		log.Printf("Error encoding schedule response: %v", err)
	}

	log.Printf("Successfully returned schedule %d with %d tasks", id, len(schedule.Tasks))
}

// @Summary Start a visit
// @Description Logs the start time and geolocation for a specific schedule. This marks the beginning of a caregiver's visit.
// @Tags schedules
// @Accept  json
// @Produce  json
// @Param id path int true "Schedule ID" example(1)
// @Param visit body models.StartVisitRequest true "Start Visit Payload"
// @Success 200 {object} map[string]string "Visit started successfully"
// @Failure 400 {object} map[string]string "Invalid request body or Schedule ID"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /schedules/{id}/start [post]
func (h *Handler) handleStartVisit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	log.Printf("POST /schedules/%s/start - Request received", idStr)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid schedule ID: %s", idStr)
		http.Error(w, "Invalid schedule ID", http.StatusBadRequest)
		return
	}

	var req models.StartVisitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding start visit request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Timestamp.IsZero() {
		req.Timestamp = time.Now()
	}

	if err := h.store.StartVisit(id, req); err != nil {
		log.Printf("Error starting visit for schedule %d: %v", id, err)
		http.Error(w, "Failed to start visit", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message":    "Visit started successfully",
		"scheduleId": idStr,
		"timestamp":  req.Timestamp.Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding start visit response: %v", err)
	}

	log.Printf("Successfully started visit for schedule %d", id)
}

// @Summary End a visit
// @Description Logs the end time and geolocation for a specific schedule. This marks the completion of a caregiver's visit.
// @Tags schedules
// @Accept  json
// @Produce  json
// @Param id path int true "Schedule ID" example(1)
// @Param visit body models.EndVisitRequest true "End Visit Payload"
// @Success 200 {object} map[string]string "Visit ended successfully"
// @Failure 400 {object} map[string]string "Invalid request body or Schedule ID"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /schedules/{id}/end [post]
func (h *Handler) handleEndVisit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	log.Printf("POST /schedules/%s/end - Request received", idStr)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("Invalid schedule ID: %s", idStr)
		http.Error(w, "Invalid schedule ID", http.StatusBadRequest)
		return
	}

	var req models.EndVisitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding end visit request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Timestamp.IsZero() {
		req.Timestamp = time.Now()
	}

	if err := h.store.EndVisit(id, req); err != nil {
		log.Printf("Error ending visit for schedule %d: %v", id, err)
		http.Error(w, "Failed to end visit", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message":    "Visit ended successfully",
		"scheduleId": idStr,
		"timestamp":  req.Timestamp.Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding end visit response: %v", err)
	}

	log.Printf("Successfully ended visit for schedule %d", id)
}

// @Summary Update task status
// @Description Updates the status of a specific task. The status can be 'pending', 'completed', or 'not_completed'. If the status is 'not_completed', a reason is required.
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param taskId path int true "Task ID" example(1)
// @Param task body models.UpdateTaskRequest true "Update Task Payload"
// @Success 200 {object} map[string]string "Task updated successfully"
// @Failure 400 {object} map[string]string "Invalid request body or Task ID"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /tasks/{taskId}/update [post]
func (h *Handler) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskIdStr := vars["taskId"]

	log.Printf("POST /tasks/%s/update - Request received", taskIdStr)

	taskId, err := strconv.Atoi(taskIdStr)
	if err != nil {
		log.Printf("Invalid task ID: %s", taskIdStr)
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding update task request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Status != "completed" && req.Status != "not_completed" && req.Status != "pending" {
		log.Printf("Invalid task status: %s", req.Status)
		http.Error(w, "Invalid task status. Must be 'completed', 'not_completed', or 'pending'", http.StatusBadRequest)
		return
	}

	if err := h.store.UpdateTaskStatus(taskId, req); err != nil {
		log.Printf("Error updating task %d: %v", taskId, err)
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "Task updated successfully",
		"taskId":  taskIdStr,
		"status":  req.Status,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding update task response: %v", err)
	}

	log.Printf("Successfully updated task %d to status %s", taskId, req.Status)
}

// @Summary Get dashboard stats
// @Description Retrieves dashboard statistics, including counts for total, missed, upcoming, and completed schedules for the current day.
// @Tags stats
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Stats "Dashboard statistics"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /stats [get]
func (h *Handler) handleGetStats(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET /stats - Request received")

	stats, err := h.store.GetStats()
	if err != nil {
		log.Printf("Error fetching stats: %v", err)
		http.Error(w, "Failed to fetch stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		log.Printf("Error encoding stats response: %v", err)
	}

	log.Printf("Successfully returned stats: total=%d, missed=%d, upcoming=%d, completed=%d",
		stats.TotalSchedules, stats.MissedSchedules, stats.UpcomingToday, stats.CompletedToday)
}
