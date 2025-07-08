package db

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func InitDB() {
	var err error
	var connectString = os.Getenv("DB_CONN")
	DB, err = sqlx.Connect("postgres", connectString)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
}
