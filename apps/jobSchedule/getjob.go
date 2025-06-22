package jobschedule

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"jobservice/dbconnection"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
)

type Scheduler struct {
	db     *sql.DB
	mu     sync.Mutex
	jobs   map[string]*Job
	ticker *time.Ticker
	quit   chan struct{}
}

func NewScheduler(db *sql.DB) *Scheduler {
	return &Scheduler{
		db:     db,
		jobs:   make(map[string]*Job),
		ticker: time.NewTicker(1 * time.Minute),
		quit:   make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	go func() {
		for {
			select {
			case now := <-s.ticker.C:
				s.checkAndRun(now)
			case <-s.quit:
				return
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	close(s.quit)
	s.ticker.Stop()
}

func (s *Scheduler) checkAndRun(now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, job := range s.jobs {
		if shouldRun(job.Schedule, now) && now.After(job.NextRun) {
			go executeJob(job)
			job.LastRun = now
			job.NextRun = nextRun(job.Schedule, now)
			s.db.Exec("UPDATE jobs SET last_run=?, next_run=? WHERE id=?", job.LastRun, job.NextRun, job.ID)
		}
	}
}

func shouldRun(expr string, now time.Time) bool {
	parts := strings.Split(expr, " ")
	if len(parts) != 5 {
		return false
	}
	min, _ := strconv.Atoi(parts[0])
	hour, _ := strconv.Atoi(parts[1])
	dow, _ := strconv.Atoi(parts[4])
	return now.Minute() == min && now.Hour() == hour && int(now.Weekday()) == dow
}

func nextRun(expr string, from time.Time) time.Time {
	for i := 1; i < 1440; i++ {
		t := from.Add(time.Duration(i) * time.Minute)
		if shouldRun(expr, t) {
			return t
		}
	}
	return from.Add(24 * time.Hour)
}
func executeJob(job *Job) {
	fmt.Printf("Running job %s at %s\n", job.Name, time.Now())
	// dummy work
}

// Job structure
type Job struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Schedule  string    `json:"schedule"` // cron or interval expression
	LastRun   time.Time `json:"last_run"`
	NextRun   time.Time `json:"next_run"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	Jobs      = make(map[string]*Job)
	jobsMutex sync.RWMutex
	scheduler = gocron.NewScheduler(time.UTC)
)

// =============================================================================
// Fetch all scheduled jobs
// Purpose: Retrieve the full list of jobs currently scheduled
// Request: GET /jobs
// Response: 200 OK, JSON { status: "S", data: [ Job, ... ] }
//
//	405 on wrong method, 500 on internal error
//
// Author: Rajkumar
// Date: 22 Jun 2025
// =============================================================================
func GetJobs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	log.Println("GetJobs:", r.Method)

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, `{"status":"E","error":"Method not allowed"}`)
		return
	}

	jobsMutex.RLock()

	lListLobs, lErr := GetAllJobs()
	if lErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status":"E","error":"Internal error"}`)
		return
	}

	jobsMutex.RUnlock()

	resp := map[string]interface{}{"status": "S", "data": lListLobs}
	b, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status":"E","error":"Internal error"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", b)
}

func GetAllJobs() ([]Job, error) {
	rows, err := dbconnection.DB.Query(
		`SELECT id, name, schedule, last_run, next_run, created_at FROM jobs`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		var j Job
		if err := rows.Scan(&j.ID, &j.Name, &j.Schedule,
			&j.LastRun, &j.NextRun, &j.CreatedAt); err != nil {
			return nil, err
		}
		jobs = append(jobs, j)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return jobs, nil
}
