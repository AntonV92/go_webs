package main

import (
	"database/sql"
	"time"
)

type User struct {
	id           int
	name         string
	password     string
	token        sql.NullString
	token_update time.Time
}

const (
	tokenExpiredMinutes = 60
)

func main() {

}

func login(login string, pass string) string {
	db := getDbConnection()

	user := User{}

	userRecord := db.QueryRow("SELECT * FROM users WHERE name = $1", login)
	userRecord.Scan(&user.id, &user.name, &user.password, &user.token, &user.token_update)

	if !checkPassword(pass, user.password) {
		return "Wrong password"
	}

	lastUpdated := user.token_update
	lastUpdated = time.Date(
		lastUpdated.Year(), lastUpdated.Month(), lastUpdated.Day(),
		lastUpdated.Hour(), lastUpdated.Minute(), 0, 0, time.Local)

	tokenCreated := time.Now().Sub(lastUpdated).Minutes()

	var token string
	if user.token.String == "" || int(tokenCreated) > tokenExpiredMinutes {
		genToken := generateToken(30)
		updatedDate := time.Now().Format(time.DateTime)
		db.QueryRow("UPDATE users SET token = $1, token_update = $2 WHERE id = $3 returning token;",
			genToken, updatedDate, user.id).Scan(&token)
	} else {
		db.QueryRow("SELECT token FROM users WHERE id = $1", user.id).Scan(&token)
	}

	return token
}
