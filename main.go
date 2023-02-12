package main

import "fmt"

type User struct {
	id           int
	name         string
	password     string
	token        string
	token_update string
}

func main() {
	fmt.Println(generateToken(40))
}
