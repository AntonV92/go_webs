package handlers

import (
	"authorization/user"
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

	user, isLoggedUser := user.UsersStorage[userId]
	if !isLoggedUser {
		conn.WriteMessage(websocket.TextMessage, []byte("Login session is expired"))
		return
	}

	ClientsOnline.ClientsList[user.Id] = user.Name
	user.WsConn = *conn
	ClientsEvents <- true

	fmt.Printf("Connected: %s\n", user.Name)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
	}
}
