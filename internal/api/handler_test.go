package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Ammar022/evv-logger-backend/internal/models"
	"github.com/gorilla/mux"
)

type MockStore struct {
	schedules []models.Schedule
	stats     *models.Stats
}

func (m *MockStore) Init() error                                                        { return nil }
func (m *MockStore) SeedData() error                                                   { return nil }
func (m *MockStore) GetSchedules() ([]models.Schedule, error)                          { return m.schedules, nil }
func (m *MockStore) GetTodaySchedules() ([]models.Schedule, error)                     { return m.schedules, nil }
func (m *MockStore) GetScheduleByID(id int) (*models.Schedule, error)                  { return &m.schedules[0], nil }
func (m *MockStore) GetStats() (*models.Stats, error)                                  { return m.stats, nil }
func (m *MockStore) StartVisit(scheduleID int, req models.StartVisitRequest) error     { return nil }
func (m *MockStore) EndVisit(scheduleID int, req models.EndVisitRequest) error         { return nil }
func (m *MockStore) UpdateTaskStatus(taskID int, req models.UpdateTaskRequest) error   { return nil }
func (m *MockStore) GetTasksByScheduleID(scheduleID int) ([]models.Task, error) {
	return []models.Task{
		{ID: 1, ScheduleID: scheduleID, Description: "Test task", Status: "pending"},
	}, nil
}

func TestGetSchedules(t *testing.T) {
	mockStore := &MockStore{
		schedules: []models.Schedule{
			{
				ID:         1,
				ClientName: "John Doe",
				ShiftStart: time.Now(),
				ShiftEnd:   time.Now().Add(8 * time.Hour),
				Location:   "123 Main St",
				VisitStatus: "upcoming",
			},
		},
	}

	handler := NewHandler(mockStore)

	req, err := http.NewRequest("GET", "/schedules", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.handleGetSchedules(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var schedules []models.Schedule
	if err := json.Unmarshal(rr.Body.Bytes(), &schedules); err != nil {
		t.Errorf("Could not unmarshal response: %v", err)
	}

	if len(schedules) != 1 {
		t.Errorf("Expected 1 schedule, got %d", len(schedules))
	}

	if schedules[0].ClientName != "John Doe" {
		t.Errorf("Expected client name 'John Doe', got '%s'", schedules[0].ClientName)
	}
}

func TestGetStats(t *testing.T) {
	mockStore := &MockStore{
		stats: &models.Stats{
			TotalSchedules:  10,
			MissedSchedules: 2,
			UpcomingToday:   3,
			CompletedToday:  5,
		},
	}

	handler := NewHandler(mockStore)

	req, err := http.NewRequest("GET", "/stats", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	handler.handleGetStats(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var stats models.Stats
	if err := json.Unmarshal(rr.Body.Bytes(), &stats); err != nil {
		t.Errorf("Could not unmarshal response: %v", err)
	}

	if stats.TotalSchedules != 10 {
		t.Errorf("Expected total schedules 10, got %d", stats.TotalSchedules)
	}
}

func TestStartVisit(t *testing.T) {
	mockStore := &MockStore{}
	handler := NewHandler(mockStore)

	visitReq := models.StartVisitRequest{
		Timestamp: time.Now(),
		Location: models.Location{
			Latitude:  40.7128,
			Longitude: -74.0060,
		},
	}

	body, _ := json.Marshal(visitReq)
	req, err := http.NewRequest("POST", "/schedules/1/start", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	rr := httptest.NewRecorder()
	handler.handleStartVisit(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("Could not unmarshal response: %v", err)
	}

	if response["message"] != "Visit started successfully" {
		t.Errorf("Expected success message, got '%s'", response["message"])
	}
}
