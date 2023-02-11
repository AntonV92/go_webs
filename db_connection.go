package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func getDbConnection() sql.DB {

	envErr := godotenv.Load()

	if envErr != nil {
		log.Fatal(envErr)
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dataSource := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, pass, dbName)

	db, err := sql.Open("postgres", dataSource)
	if err != nil {
		fmt.Printf("DB error: %s", err)
	}

	return *db
}
