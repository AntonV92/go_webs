package handlers

import (
	"authorization/user"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

type WsLoginData struct {
	UserId int    `json:"user_id"`
	Token  string `json:"token"`
}

type OnlineClients struct {
	ClientsList map[int]string `json:"online_clients"`
}

type ClientMessage struct {
	ToUserId    int    `json:"to_user_id"`
	FromUserId  int    `json:"from_user_id"`
	MessageText string `json:"message_text"`
}

var ClientsOnline = OnlineClients{
	ClientsList: make(map[int]string),
}
var ClientsEvents = make(chan bool)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {

		userId := r.URL.Query().Get("user_id")
		token := r.URL.Query().Get("token")

		return user.CheckToken(userId, token)
	},
}

func WsHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	userId, strConvErr := strconv.Atoi(r.URL.Query().Get("user_id"))

	defer func() {
		delete(ClientsOnline.ClientsList, userId)
		ClientsEvents <- true
	}()

	if strConvErr != nil {
		fmt.Println(strConvErr)
		return
	}

	loggedUser, isLoggedUser := user.UsersStorage[userId]
	if !isLoggedUser {
		conn.WriteMessage(websocket.TextMessage, []byte("Login session is expired"))
		return
	}

	ClientsOnline.ClientsList[loggedUser.Id] = loggedUser.Name
	loggedUser.WsConn = *conn
	ClientsEvents <- true

	fmt.Printf("Connected: %s\n", loggedUser.Name)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read ws message error: ", err)
			break
		}

		sendMessageErr := sendClientMessage(message)
		if sendMessageErr != nil {
			fmt.Printf("send client message error: %v\n", sendMessageErr)
			continue
		}
	}
}

func sendClientMessage(message []byte) error {

	mes := ClientMessage{}
	decodeError := json.Unmarshal(message, &mes)
	if decodeError != nil {
		return decodeError
	}
	sendError := user.UsersStorage[mes.ToUserId].WsConn.WriteMessage(websocket.TextMessage, []byte(message))
	if sendError != nil {
		return sendError
	}

	return nil
}
