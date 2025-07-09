package infrastructure

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func InitDB() *sqlx.DB {
	var connectString = os.Getenv("DB_CONN")
	db, err := sqlx.Connect("postgres", connectString)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	return db
}
