package main

import (
	"log"
	"math/rand"

	"golang.org/x/crypto/bcrypt"
)

const (
	chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

// get hashed string from given password
func passwordHash(pass string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), 0)
	if err != nil {
		log.Fatal(err)
	}

	return string(hash)
}

// compare password with string hash
func checkPassword(pass string, passwordHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(pass))
	if err != nil {
		return false
	}

	return true
}

func generateToken(lenght int) string {

	token := ""

	for i := 0; i < lenght; i++ {
		randIndex := rand.Intn(len(chars) - 1)
		token += string(chars[randIndex])
	}

	return token
}
