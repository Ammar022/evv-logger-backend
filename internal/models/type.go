package models

import "time"

// Location represents geographic coordinates
// swagger:model
// @Example {"latitude":40.7128,"longitude":-74.0060}
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Schedule struct {
	ID         int       `json:"id"`
	ClientName string    `json:"clientName"`
	ShiftStart time.Time `json:"shiftStart"`
	ShiftEnd   time.Time `json:"shiftEnd"`
	Location   string    `json:"location"`
	VisitStatus string   `json:"visitStatus"`
	Tasks      []Task    `json:"tasks,omitempty"`
}

type Task struct {
	ID          int    `json:"id"`
	ScheduleID  int    `json:"scheduleId"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Reason      string `json:"reason,omitempty"`
}

type VisitLog struct {
	ID         int       `json:"id"`
	ScheduleID int       `json:"scheduleId"`
	LogType    string    `json:"logType"`
	Timestamp  time.Time `json:"timestamp"`
	Location   Location  `json:"location"`
}

type Stats struct {
	TotalSchedules    int `json:"totalSchedules"`
	MissedSchedules   int `json:"missedSchedules"`
	UpcomingToday     int `json:"upcomingToday"`
	CompletedToday    int `json:"completedToday"`
}

type StartVisitRequest struct {
	Timestamp time.Time `json:"timestamp" example:"2023-05-15T09:00:00Z"`
	Location  Location  `json:"location"`
}

type EndVisitRequest struct {
	Timestamp time.Time `json:"timestamp" example:"2023-05-15T17:00:00Z"`
	Location  Location  `json:"location"`
}

type UpdateTaskRequest struct {
	Status string `json:"status" example:"completed"`
	Reason string `json:"reason,omitempty" example:"reason for update"`
}
