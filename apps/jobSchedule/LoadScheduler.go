package jobschedule

import (
	"fmt"
	"jobservice/dbconnection"
	"strconv"
	"strings"
	"time"
)

func LoadJobsFromDB() (map[string]*Job, error) {
	rows, err := dbconnection.DB.Query(`
        SELECT id, name, schedule, last_run, next_run, created_at
        FROM jobs
    `)
	if err != nil {
		return nil, fmt.Errorf("query jobs: %w", err)
	}
	defer rows.Close()

	jobs := make(map[string]*Job)
	for rows.Next() {
		var j Job
		if err := rows.Scan(
			&j.ID, &j.Name, &j.Schedule,
			&j.LastRun, &j.NextRun, &j.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan job: %w", err)
		}

		// If NextRun is zero, compute it
		if j.NextRun.IsZero() {
			j.NextRun = calculateNextRun(j.Schedule, time.Now())
		}

		jobs[j.ID] = &j
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows err: %w", err)
	}
	return jobs, nil
}

func calculateNextRun(expr string, from time.Time) time.Time {
	parts := strings.Split(expr, " ")
	if len(parts) != 5 {
		return from.Add(time.Hour) // fallback
	}
	min, _ := strconv.Atoi(parts[0])
	hour, _ := strconv.Atoi(parts[1])
	dow, _ := strconv.Atoi(parts[4])

	for i := 0; i < 1440; i++ {
		t := from.Add(time.Duration(i) * time.Minute)
		if t.Minute() == min && t.Hour() == hour && int(t.Weekday()) == dow {
			return t
		}
	}
	return from.Add(24 * time.Hour)
}
