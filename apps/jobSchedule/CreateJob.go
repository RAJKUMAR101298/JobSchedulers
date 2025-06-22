package jobschedule

import (
	"encoding/json"
	"fmt"
	"jobservice/dbconnection"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// =============================================================================
// Create a new scheduled job
// Purpose: Add a new job to the scheduler by providing a name and schedule
// Request: POST /jobs
//
//	JSON body { "name": "<job name>", "schedule": "<HH:MM>" }
//
// Response: 200 OK, JSON { status: "S", data: Job }
//
//	400 on invalid body, 405 on wrong method, 500 on internal error
//
// Author: Rajkumar
// Date: 22 Jun 2025
// =============================================================================
func CreateJob(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Println("CreateJob:", r.Method)

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, `{"status":"E","error":"Method not allowed"}`)
		return
	}

	var in struct {
		Name     string `json:"name"`
		Schedule string `json:"schedule"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"status":"E","error":"Invalid payload"}`)
		return
	}
	defer r.Body.Close()

	job := &Job{
		ID:        uuid.New().String(),
		Name:      in.Name,
		Schedule:  in.Schedule,
		CreatedAt: time.Now(),
	}
	// scheduleJob(job)

	lErr := InsertJob(job)
	if lErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status":"E","error":"Internal error"}`)
		return
	}

	resp := map[string]interface{}{"status": "S", "data": job}
	b, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status":"E","error":"Internal error"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", b)
}

func InsertJob(job *Job) error {
	_, err := dbconnection.DB.Exec(`
        INSERT INTO jobs (id, name, schedule, last_run, next_run, created_at)
        VALUES (?, ?, ?, ?, ?, ?)`,
		job.ID, job.Name, job.Schedule,
		job.LastRun, job.NextRun, job.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("InsertJob exec error: %w", err)
	}
	return nil
}
