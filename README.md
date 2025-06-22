
# 📆 Scheduler Microservice

A Go-based microservice for creating, storing, and executing scheduled jobs with MySQL persistence.

---

## 🏗 Project Structure

```text
.
├── cmd/           # Application entrypoint: DB connect, scheduler, and HTTP setup
├── db/            # MySQL connection logic
├── jobschedule/   # Job struct, scheduler, and HTTP handlers
├── go.mod
└── go.sum


🛠 Prerequisites
Go 1.20+

MySQL 5.7+

MySQL client & Go CLI

📦 Setup & Installation
1. Initialize Go Module
mkdir scheduler-service && cd scheduler-service
go mod init github.com/yourusername/scheduler-service
go mod tidy


2. Install MySQL & Create Schema
sudo apt install mysql-server

CREATE DATABASE schedulerdb;
USE schedulerdb;

CREATE TABLE IF NOT EXISTS jobs (
  id VARCHAR(36) PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  schedule VARCHAR(100) NOT NULL,
  last_run DATETIME NULL,
  next_run DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_next_run ON jobs(next_run);

| Endpoint           | Request Body                                       | Success Response                                                 | Error Responses                                                                                       |
| ------------------ | -------------------------------------------------- | ---------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------- |
| **GET /jobs**      | –                                                  | `{"status":"S","data":[{…job…}, …]}`                             | –                                                                                                     |
| **GET /jobs/{id}** | –                                                  | `{"status":"S","data":{…job…}}`                                  | `404: {"status":"E","error":"Job not found"}` <br> `405: {"status":"E","error":"Method not allowed"}` |
| **POST /jobs**     | `{"name": "ExampleJob", "schedule": "0 12 * * *"}` | `{"status":"S","data":{…new job with id,next_run, created_at…}}` | `400: invalid payload` <br> `405: method not allowed` <br> `500: DB insert failed`                    |


👤 Author
Rajkumar • June 22, 2025