package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func ConnectDatabase() (*sql.DB, error) {
	host := "localhost"
	port := 5432
	user := "postgres"
	dbname := "postgres"
	pass := "haseem123"
	psqlSetup := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable", host, port, user, dbname, pass)
	var err error
	db, err = sql.Open("postgres", psqlSetup)
	if err != nil {
		fmt.Println("There is an error while connecting to the database ", err)
		fmt.Println(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Could not establish a connection to the database: %v", err)
	}
	fmt.Println("Successfully connected to the database")
	return db, nil
}
