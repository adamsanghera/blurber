package main

import (
	"log"
	"net/http"

	"./handler"
)

func main() {
	log.Printf("SERVER: Starting...")

	http.HandleFunc("/login/", handler.Auth)
	http.HandleFunc("/feed/", handler.Feed)
	http.HandleFunc("/blurb/", handler.Blurb)
	http.HandleFunc("/subscribe/", handler.Subscribe)
	http.HandleFunc("/unsubscribe/", handler.Unsubscribe)
	http.HandleFunc("/profile/", handler.Profile)
	http.HandleFunc("/removeAcc/", handler.RemoveAcc)

	http.ListenAndServe(":4000", nil)
}
