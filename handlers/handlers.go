package handlers

import (
	"authorization/user"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LoginForm struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Read body error")
	}

	var loginData LoginForm

	jsonErr := json.Unmarshal(body, &loginData)
	if jsonErr != nil {
		w.Write([]byte("json login data error"))
		fmt.Println(jsonErr)
	}
	user, loginErr := user.Login(loginData.Login, loginData.Password)

	if loginErr != nil {
		w.Write([]byte("Login user error"))
		fmt.Println(loginErr)
		return
	}

	rData := WsLoginData{
		UserId: user.Id,
		Token:  user.Token.String,
	}

	response, responseJsonError := json.Marshal(rData)

	if responseJsonError != nil {
		fmt.Printf("response json error %v", responseJsonError)
	}

	w.Write([]byte(response))

}
