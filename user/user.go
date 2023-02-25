package user

import (
	"authorization/db"
	"authorization/passwd"
	"database/sql"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

type User struct {
	Id           int            `json:"id"`
	Name         string         `json:"name"`
	Password     string         `json:"password"`
	Token        sql.NullString `json:"token"`
	Token_update time.Time      `json:"token_update"`
	WsConn       websocket.Conn
}

var UsersStorage = make(map[int]*User)

func Login(login string, pass string) (User, error) {
	db := db.DbConn

	user := User{}

	userRecord := db.QueryRow("SELECT * FROM users WHERE name = $1 LIMIT 1;", login)
	userRecord.Scan(&user.Id, &user.Name, &user.Password, &user.Token, &user.Token_update)

	if !passwd.CheckPassword(pass, user.Password) {
		return User{}, fmt.Errorf("Wrong password")
	}

	lastUpdated := user.Token_update
	lastUpdated = time.Date(
		lastUpdated.Year(), lastUpdated.Month(), lastUpdated.Day(),
		lastUpdated.Hour(), lastUpdated.Minute(), 0, 0, time.Local)

	tokenCreated := time.Now().Sub(lastUpdated).Minutes()

	var token sql.NullString
	if user.Token.String == "" || int(tokenCreated) > passwd.TokenExpiredMinutes {
		genToken := passwd.GenerateToken(30)
		updatedDate := time.Now().Format(time.DateTime)
		db.QueryRow("UPDATE users SET token = $1, token_update = $2 WHERE id = $3 returning token;",
			genToken, updatedDate, user.Id).Scan(&token)
	} else {
		db.QueryRow("SELECT token FROM users WHERE id = $1 LIMIT 1;", user.Id).Scan(&token)
	}

	user.Token = token

	UsersStorage[user.Id] = &user

	return user, nil
}

func CheckToken(userId string, token string) bool {
	db := db.DbConn
	user := User{}
	record := db.QueryRow("SELECT id FROM users WHERE id = $1 AND token = $2 LIMIT 1;", userId, token)

	err := record.Scan(&user.Id)
	if err != nil {
		fmt.Println("Scan user error")
		return false
	}

	return true
}
