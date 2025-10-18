package repo

import (
	"database/sql"

	"github.com/Ammar022/evv-logger-backend/internal/models"
)

// Store defines all database operations
type Store interface {
	Init() error
	SeedData() error
	GetSchedules() ([]models.Schedule, error)
	GetTodaySchedules() ([]models.Schedule, error)
	GetScheduleByID(id int) (*models.Schedule, error)
	GetStats() (*models.Stats, error)
	StartVisit(scheduleID int, req models.StartVisitRequest) error
	EndVisit(scheduleID int, req models.EndVisitRequest) error
	UpdateTaskStatus(taskID int, req models.UpdateTaskRequest) error
	GetTasksByScheduleID(scheduleID int) ([]models.Task, error)
}

// PostgresStore implements the Store interface
type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func (s *PostgresStore) Init() error {
	createTablesSQL := `
    CREATE TABLE IF NOT EXISTS schedules (
        id SERIAL PRIMARY KEY,
        client_name VARCHAR(255) NOT NULL,
        shift_start TIMESTAMP NOT NULL,
        shift_end TIMESTAMP NOT NULL,
        location VARCHAR(255),
        visit_status VARCHAR(50) DEFAULT 'upcoming',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS tasks (
        id SERIAL PRIMARY KEY,
        schedule_id INTEGER REFERENCES schedules(id) ON DELETE CASCADE,
        description VARCHAR(500) NOT NULL,
        status VARCHAR(50) DEFAULT 'pending',
        reason TEXT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS visit_logs (
        id SERIAL PRIMARY KEY,
        schedule_id INTEGER REFERENCES schedules(id) ON DELETE CASCADE,
        log_type VARCHAR(10) NOT NULL CHECK (log_type IN ('start', 'end')),
        timestamp TIMESTAMP NOT NULL,
        latitude DECIMAL(10, 8),
        longitude DECIMAL(11, 8),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    `
	_, err := s.db.Exec(createTablesSQL)
	return err
}

// GetSchedules retrieves all schedules
func (s *PostgresStore) GetSchedules() ([]models.Schedule, error) {
	query := `
		SELECT id, client_name, shift_start, shift_end, location, visit_status 
		FROM schedules 
		ORDER BY shift_start ASC
	`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.Schedule
	for rows.Next() {
		var schedule models.Schedule
		err := rows.Scan(
			&schedule.ID,
			&schedule.ClientName,
			&schedule.ShiftStart,
			&schedule.ShiftEnd,
			&schedule.Location,
			&schedule.VisitStatus,
		)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}
	return schedules, nil
}

// GetTodaySchedules retrieves schedules for today
func (s *PostgresStore) GetTodaySchedules() ([]models.Schedule, error) {
	query := `
		SELECT id, client_name, shift_start, shift_end, location, visit_status 
		FROM schedules 
		WHERE DATE(shift_start) = CURRENT_DATE
		ORDER BY shift_start ASC
	`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.Schedule
	for rows.Next() {
		var schedule models.Schedule
		err := rows.Scan(
			&schedule.ID,
			&schedule.ClientName,
			&schedule.ShiftStart,
			&schedule.ShiftEnd,
			&schedule.Location,
			&schedule.VisitStatus,
		)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}
	return schedules, nil
}

// GetScheduleByID retrieves a schedule by ID with its tasks
func (s *PostgresStore) GetScheduleByID(id int) (*models.Schedule, error) {
	query := `
		SELECT id, client_name, shift_start, shift_end, location, visit_status 
		FROM schedules 
		WHERE id = $1
	`
	var schedule models.Schedule
	err := s.db.QueryRow(query, id).Scan(
		&schedule.ID,
		&schedule.ClientName,
		&schedule.ShiftStart,
		&schedule.ShiftEnd,
		&schedule.Location,
		&schedule.VisitStatus,
	)
	if err != nil {
		return nil, err
	}

	// Get associated tasks
	tasks, err := s.GetTasksByScheduleID(id)
	if err != nil {
		return nil, err
	}
	schedule.Tasks = tasks

	return &schedule, nil
}

// GetTasksByScheduleID retrieves tasks for a specific schedule
func (s *PostgresStore) GetTasksByScheduleID(scheduleID int) ([]models.Task, error) {
	query := `
		SELECT id, schedule_id, description, status, COALESCE(reason, '') 
		FROM tasks 
		WHERE schedule_id = $1
		ORDER BY id ASC
	`
	rows, err := s.db.Query(query, scheduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID,
			&task.ScheduleID,
			&task.Description,
			&task.Status,
			&task.Reason,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// GetStats retrieves dashboard statistics
func (s *PostgresStore) GetStats() (*models.Stats, error) {
	var stats models.Stats
	
	// Total schedules
	err := s.db.QueryRow("SELECT COUNT(*) FROM schedules").Scan(&stats.TotalSchedules)
	if err != nil {
		return nil, err
	}
	
	// Missed schedules (today)
	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM schedules 
		WHERE DATE(shift_start) = CURRENT_DATE AND visit_status = 'missed'
	`).Scan(&stats.MissedSchedules)
	if err != nil {
		return nil, err
	}
	
	// Upcoming today
	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM schedules 
		WHERE DATE(shift_start) = CURRENT_DATE AND visit_status = 'upcoming'
	`).Scan(&stats.UpcomingToday)
	if err != nil {
		return nil, err
	}
	
	// Completed today
	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM schedules 
		WHERE DATE(shift_start) = CURRENT_DATE AND visit_status = 'completed'
	`).Scan(&stats.CompletedToday)
	if err != nil {
		return nil, err
	}
	
	return &stats, nil
}

// StartVisit logs the start of a visit
func (s *PostgresStore) StartVisit(scheduleID int, req models.StartVisitRequest) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	// Insert visit log
	_, err = tx.Exec(`
		INSERT INTO visit_logs (schedule_id, log_type, timestamp, latitude, longitude)
		VALUES ($1, 'start', $2, $3, $4)
	`, scheduleID, req.Timestamp, req.Location.Latitude, req.Location.Longitude)
	if err != nil {
		return err
	}
	
	// Update schedule status to 'in_progress' or similar
	_, err = tx.Exec(`
		UPDATE schedules SET visit_status = 'in_progress' WHERE id = $1
	`, scheduleID)
	if err != nil {
		return err
	}
	
	return tx.Commit()
}

// EndVisit logs the end of a visit
func (s *PostgresStore) EndVisit(scheduleID int, req models.EndVisitRequest) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	// Insert visit log
	_, err = tx.Exec(`
		INSERT INTO visit_logs (schedule_id, log_type, timestamp, latitude, longitude)
		VALUES ($1, 'end', $2, $3, $4)
	`, scheduleID, req.Timestamp, req.Location.Latitude, req.Location.Longitude)
	if err != nil {
		return err
	}
	
	// Update schedule status to 'completed'
	_, err = tx.Exec(`
		UPDATE schedules SET visit_status = 'completed' WHERE id = $1
	`, scheduleID)
	if err != nil {
		return err
	}
	
	return tx.Commit()
}

// UpdateTaskStatus updates the status of a task
func (s *PostgresStore) UpdateTaskStatus(taskID int, req models.UpdateTaskRequest) error {
	_, err := s.db.Exec(`
		UPDATE tasks SET status = $1, reason = $2 WHERE id = $3
	`, req.Status, req.Reason, taskID)
	return err
}
