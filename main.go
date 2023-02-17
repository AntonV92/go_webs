package main

import (
	"encoding/json"
	"fmt"
	"io"
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

	http.HandleFunc("/login", handleLogin)

	log.Fatal(http.ListenAndServe(":8000", nil))
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var loginData LoginForm

	jsonErr := json.Unmarshal(body, &loginData)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	user, loginErr := login(loginData.Login, loginData.Password)

	if loginErr != nil {
		fmt.Println(loginErr)
		return
	}

	w.Write([]byte(user.token.String))

	fmt.Println(user.token.String)
}
