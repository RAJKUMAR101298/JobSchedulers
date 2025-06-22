package jobschedule

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"jobservice/dbconnection"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// =============================================================================
// Fetch a scheduled job by ID
// Purpose: Retrieve details of a single job using its ID
// Request: GET /jobs/{id}
// Response: 200 OK, JSON { status: "S", data: Job }
//
//	404 if not found, 405 on wrong method, 500 on internal error
//
// Author: Rajkumar
// Date: 22 Jun 2025
// =============================================================================
func GetJobByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Println("GetJobByID:", r.Method)

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, `{"status":"E","error":"Method not allowed"}`)
		return
	}

	id := mux.Vars(r)["id"]
	jobsMutex.RLock()
	lJob, lErr := GetJobById(id)
	if lErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status":"E","error":"Internal error"}`)
		return
	}

	jobsMutex.RUnlock()

	resp := map[string]interface{}{"status": "S", "data": lJob}
	b, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"status":"E","error":"Internal error"}`)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", b)
}

func GetJobById(id string) (Job, error) {
	var lJobById Job
	err := dbconnection.DB.QueryRow(
		`SELECT id, name, schedule, last_run, next_run, created_at FROM jobs WHERE id = ?`, id,
	).Scan(&lJobById.ID, &lJobById.Name, &lJobById.Schedule, &lJobById.LastRun, &lJobById.NextRun, &lJobById.CreatedAt)
	if err == sql.ErrNoRows {
		return lJobById, nil
	}
	if err != nil {
		return lJobById, err
	}
	return lJobById, nil
}
