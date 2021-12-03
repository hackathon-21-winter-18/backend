package main

import "log"

func main() {
	log.Println("Server Started.")

	db, err := model.InitDB()
}
