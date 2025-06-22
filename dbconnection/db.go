package dbconnection

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB // Global DB variable

func ConnectDB() error {
	// Predefined connection values
	user := "root"
	password := "root"
	host := "localhost"
	port := "3306"
	dbname := "schedulerdb"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("failed to open DB: %v", err)
	}

	if err := DB.Ping(); err != nil {
		log.Println(err)

		return fmt.Errorf("failed to connect to DB: %v", err)
	}

	return nil
}
