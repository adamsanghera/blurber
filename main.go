package main

import (
	"log"
	"net/http"

	"./handler"
)

func main() {
	log.Printf("SERVER: Starting...")

	http.HandleFunc("/login/", handler.Login)
	http.HandleFunc("/feed/", handler.Feed)
	http.HandleFunc("/blurb/", handler.Blurb)

	http.ListenAndServe(":4000", nil)
}
