package main

import (
	"html/template"
	"log"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("./static/login/index.html")
		t.Execute(w, nil)
	}
}

func main() {
	log.Printf("SERVER: Starting...")

	// http.Handle("/", http.FileServer(http.Dir("./static/")))

	http.HandleFunc("/", LoginHandler)

	// http.Handle("/", http.FileServer(http.Dir("./static")))
	http.ListenAndServe(":4000", nil)
}
