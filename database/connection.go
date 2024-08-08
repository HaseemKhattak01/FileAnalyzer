package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func ConnectDatabase() (*sql.DB, error) {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	pass := os.Getenv("DB_PASS")

	fmt.Printf("Connecting to database at host=%s port=%s user=%s dbname=%s\n", host, port, user, dbname)

	psqlSetup := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbname, pass)
	db, err = sql.Open("postgres", psqlSetup)
	if err != nil {
		fmt.Println("There is an error while connecting to the database ", err)
		fmt.Println(err)
		return nil, err
	}
	fmt.Println("hey im here")
	err = db.Ping()
	if err != nil {
		log.Fatalf("Could not establish a connection to the database: %v", err)
	}

	fmt.Println("Successfully connected to the database")
	return db, nil
}
