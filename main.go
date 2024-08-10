package main

import (
	"fmt"
	"log"

	"github.com/luqxus/medspace/space"
)

var (
	userBucket = "users"
)

func main() {

	user := map[string]string{
		"name": "GG",
		"age":  "28",
	}

	_ = user

	db, err := space.New()
	if err != nil {
		log.Fatal(err)
	}

	id, err := db.Insert(userBucket, user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(id)

	result, err := db.Get(userBucket, "users", "")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", result)
}
