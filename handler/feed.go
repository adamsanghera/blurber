package handler

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

func Feed(w http.ResponseWriter, r *http.Request) {
	log.Printf("HANDLERS-FEED: Request received")

	start := time.Now()
	defer func() {
		log.Printf("HANDLERS-FEED: Request serviced in %5.1f seconds", time.Since(start).Seconds())
	}()

	// Validate and renew our cookies
	validSesh, uname := validateSession(w, r)
	if !validSesh {
		return
	}

	// Build template
	t, err := template.ParseFiles("./static-assets/feed/index.html")
	if err != nil {
		panic(err)
	}

	// TODO: Implement real logic here.
	// Obtain our blurb list
	usrID, err := userDB.GetUserID(uname)
	bs := blurbDB.GetBlurbsCreatedBy(usrID)

	// Squeeze our blurbs into the template, execute
	t.Execute(w, bs)
}
