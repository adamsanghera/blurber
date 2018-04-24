package main

import (
	"log"
	"net/http"
	"os"

	"github.com/adamsanghera/blurber/handler"
)

func main() {
	log.Printf("SERVER: Preparing handlers...\n")

	http.HandleFunc("/login/", handler.Auth)
	http.HandleFunc("/feed/", handler.Feed)
	http.HandleFunc("/blurb/", handler.Blurb)
	http.HandleFunc("/subscribe/", handler.Subscribe)
	http.HandleFunc("/unsubscribe/", handler.Unsubscribe)
	http.HandleFunc("/profile/", handler.Profile)
	http.HandleFunc("/removeAcc/", handler.RemoveAcc)
	http.HandleFunc("/logout/", handler.Logout)
	http.HandleFunc("/", handler.Catcher)

	log.Printf("SERVER: Serving...")
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatalf("SERVER: Failed to listen and serve on port {%s}.  Error message: {%v}", os.Getenv("PORT"), err.Error())
	}
}
