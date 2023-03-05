package main

import (
	"authorization/db"
	"authorization/handlers"
	"authorization/user"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

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

	curr, _ := os.Getwd()

	fs := http.FileServer(http.Dir(curr + "/frontend"))

	http.Handle("/frontend/", http.StripPrefix("/frontend", fs))
	http.HandleFunc("/", handlers.HandleMain)
	http.HandleFunc("/ws", handlers.WsHandler)
	http.HandleFunc("/login", handlers.HandleLogin)

	log.Fatal(http.ListenAndServe(":8000", nil))

}

func broadcaster() {
	for {
		select {
		case <-handlers.ClientsEvents:

			for userId := range handlers.ClientsOnline.ClientsList {

				userConn := user.UsersStorage[userId].WsConn
				jsonNames, err := json.Marshal(handlers.ClientsOnline)

				if err != nil {
					fmt.Printf("Json clients event error: %v\n", err)
				}

				wserr := userConn.WriteMessage(websocket.TextMessage, jsonNames)
				if wserr != nil {
					fmt.Println(wserr)
				}
			}
		}
	}
}
