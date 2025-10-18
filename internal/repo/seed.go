package repo

import (
	"log"
	"time"
)

func (s *PostgresStore) SeedData() error {
	log.Println("Seeding database with sample data...")

	scheduleSQL := `
		INSERT INTO schedules (client_name, shift_start, shift_end, location, visit_status) 
		VALUES 
			('John Doe', $1, $2, '123 Main St, Springfield', 'upcoming'),
			('Jane Smith', $3, $4, '456 Oak Ave, Springfield', 'upcoming'),
			('Bob Johnson', $5, $6, '789 Pine Rd, Springfield', 'completed'),
			('Alice Brown', $7, $8, '321 Elm St, Springfield', 'missed')
		ON CONFLICT DO NOTHING
	`

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	
	// Today's schedules
	schedule1Start := today.Add(9 * time.Hour)  // 9 AM today
	schedule1End := today.Add(17 * time.Hour)   // 5 PM today
	schedule2Start := today.Add(10 * time.Hour) // 10 AM today
	schedule2End := today.Add(18 * time.Hour)   // 6 PM today
	
	// Yesterday's schedules
	yesterday := today.AddDate(0, 0, -1)
	schedule3Start := yesterday.Add(8 * time.Hour)  // 8 AM yesterday
	schedule3End := yesterday.Add(16 * time.Hour)   // 4 PM yesterday
	schedule4Start := yesterday.Add(14 * time.Hour) // 2 PM yesterday
	schedule4End := yesterday.Add(22 * time.Hour)   // 10 PM yesterday

	_, err := s.db.Exec(scheduleSQL, 
		schedule1Start, schedule1End,
		schedule2Start, schedule2End,
		schedule3Start, schedule3End,
		schedule4Start, schedule4End,
	)
	if err != nil {
		return err
	}

	taskSQL := `
		INSERT INTO tasks (schedule_id, description, status, reason) 
		VALUES 
			(1, 'Administer morning medication', 'pending', ''),
			(1, 'Assist with bathing', 'pending', ''),
			(1, 'Prepare lunch', 'pending', ''),
			(2, 'Check vital signs', 'pending', ''),
			(2, 'Physical therapy exercises', 'pending', ''),
			(3, 'Administer medication', 'completed', ''),
			(3, 'Assist with mobility', 'completed', ''),
			(3, 'Meal preparation', 'not_completed', 'Client was not hungry'),
			(4, 'Evening medication', 'not_completed', 'Client refused medication'),
			(4, 'Wound care', 'not_completed', 'Client was not available')
		ON CONFLICT DO NOTHING
	`

	_, err = s.db.Exec(taskSQL)
	if err != nil {
		return err
	}

	visitLogSQL := `
		INSERT INTO visit_logs (schedule_id, log_type, timestamp, latitude, longitude)
		VALUES 
			(3, 'start', $1, 39.7817, -89.6501),
			(3, 'end', $2, 39.7817, -89.6501)
		ON CONFLICT DO NOTHING
	`

	startTime := schedule3Start.Add(15 * time.Minute) // Started 15 minutes late
	endTime := schedule3End.Add(-30 * time.Minute)   // Ended 30 minutes early

	_, err = s.db.Exec(visitLogSQL, startTime, endTime)
	if err != nil {
		return err
	}

	log.Println("Sample data seeded successfully!")
	return nil
}
