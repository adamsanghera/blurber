package main

import (
	"log"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("SERVER: Handling login endpoint...\n")
	if r.Method == "GET" {
		log.Printf("SERVER: Login – GET")
		http.FileServer(http.Dir("./static/login")).ServeHTTP(w, r)
	} else {
		if r.URL.Path == "/login/existing" {

		} else if r.URL.Path == "/login/new" {

		} else {
			w.Write([]byte("Something went wrong..."))
		}
	}
}

func main() {
	log.Printf("SERVER: Starting...")

	// http.Handle("/", http.FileServer(http.Dir("./static/")))

	http.HandleFunc("/", LoginHandler)

	http.HandleFunc("/login", LoginHandler)

	// http.HandleFunc("/login", LoginHandler)

	// http.Handle("/", http.FileServer(http.Dir("./static")))
	http.ListenAndServe(":4000", nil)
}
