package main

import (
	"log"

	"golang.org/x/crypto/bcrypt"
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
