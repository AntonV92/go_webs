package main

import (
	"authorization/db"
	"authorization/handlers"
	"log"
	"net/http"
)

type LoginForm struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

const (
	tokenExpiredMinutes = 60
)

func main() {

	db.InitDbConnection()

	http.HandleFunc("/login", handlers.HandleLogin)
	http.HandleFunc("/", handlers.HandleMain)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
