package main

import (
	"log"
	"net/http"

	reg "./registration"

	"./handler"
)

func main() {
	log.Printf("SERVER: Starting...")
	increasingCounter = 0

	// Dev only - for quick loging
	userDB.Add(reg.User{"dev", "-1"}, "root")

	http.HandleFunc("/login/", handler.Login)
	http.HandleFunc("/feed/", handler.Feed)
	http.HandleFunc("/blurb/", handler.Blurb)

	http.ListenAndServe(":4000", nil)
}
