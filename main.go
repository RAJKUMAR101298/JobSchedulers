package main

import (
	"fmt"
	jobschedule "jobservice/apps/jobSchedule"
	"jobservice/dbconnection"
	"log"
	"net/http"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/gorilla/mux"
)

var (
	scheduler = gocron.NewScheduler(time.UTC)
)

func main() {
	scheduler.StartAsync()
	err := dbconnection.ConnectDB()
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	var lErr error
	jobschedule.Jobs, lErr = jobschedule.LoadJobsFromDB()
	if lErr != nil {
		log.Fatal("Job Load Error")
	}
	r := mux.NewRouter()
	r.HandleFunc("/jobs", jobschedule.GetJobs).Methods("GET")
	r.HandleFunc("/jobs/{id}", jobschedule.GetJobByID).Methods("GET")
	r.HandleFunc("/jobs", jobschedule.CreateJob).Methods("POST")

	fmt.Println("Scheduler service running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
