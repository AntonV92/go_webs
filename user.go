package main

import (
	"database/sql"
	"fmt"
	"time"
)

type User struct {
	id           int
	name         string
	password     string
	token        sql.NullString
	token_update time.Time
}

func login(login string, pass string) (User, error) {
	db := getDbConnection()

	user := User{}

	userRecord := db.QueryRow("SELECT * FROM users WHERE name = $1 LIMIT 1;", login)
	userRecord.Scan(&user.id, &user.name, &user.password, &user.token, &user.token_update)

	if !checkPassword(pass, user.password) {
		return User{}, fmt.Errorf("Wrong password")
	}

	lastUpdated := user.token_update
	lastUpdated = time.Date(
		lastUpdated.Year(), lastUpdated.Month(), lastUpdated.Day(),
		lastUpdated.Hour(), lastUpdated.Minute(), 0, 0, time.Local)

	tokenCreated := time.Now().Sub(lastUpdated).Minutes()

	var token sql.NullString
	if user.token.String == "" || int(tokenCreated) > tokenExpiredMinutes {
		genToken := generateToken(30)
		updatedDate := time.Now().Format(time.DateTime)
		db.QueryRow("UPDATE users SET token = $1, token_update = $2 WHERE id = $3 returning token;",
			genToken, updatedDate, user.id).Scan(&token)
	} else {
		db.QueryRow("SELECT token FROM users WHERE id = $1 LIMIT 1;", user.id).Scan(&token)
	}

	user.token = token

	return user, nil
}
