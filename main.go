package main

import (
	"log"
	"net/http"
	"os"

	"github.com/adamsanghera/blurber/handler"
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
	http.HandleFunc("/logout/", handler.LogOut)
	http.HandleFunc("/", handler.Catcher)

	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}
