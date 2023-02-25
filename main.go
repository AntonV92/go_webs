package main

import (
	"authorization/db"
	"authorization/handlers"
	"authorization/user"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type LoginForm struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

var Messages = make(chan string)

func main() {

	db.InitDbConnection()

	go broadcaster()
	go user.UserStorageChecker()

	http.HandleFunc("/ws", handlers.WsHandler)
	http.HandleFunc("/login", handlers.HandleLogin)

	log.Fatal(http.ListenAndServe(":8000", nil))

}

func broadcaster() {
	for {
		select {
		case <-handlers.ClientsEvents:

			for userId := range handlers.OnlineClients {

				userConn := user.UsersStorage[userId].WsConn
				jsonNames, err := json.Marshal(handlers.OnlineClients)

				if err != nil {
					fmt.Println("json error")
				}

				wserr := userConn.WriteMessage(websocket.TextMessage, jsonNames)
				if wserr != nil {
					fmt.Println(wserr)
				}
			}
		}
	}
}
