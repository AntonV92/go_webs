package handlers

import (
	"authorization/user"
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

func HandleLogin(w http.ResponseWriter, r *http.Request) {
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
	user, loginErr := user.Login(loginData.Login, loginData.Password)

	if loginErr != nil {
		fmt.Println(loginErr)
		return
	}

	w.Write([]byte(user.Token.String))

	fmt.Println(user.Token.String)
}

func HandleMain(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")

	fmt.Println(token)
}
