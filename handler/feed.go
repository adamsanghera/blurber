package handler

import (
	"html/template"
	"log"
	"net/http"
)

func Feed(w http.ResponseWriter, r *http.Request) {
	log.Printf("FEED: URL is <%s>", r.URL.Path)

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
	bs := lbl.GetBlurbsCreatedBy(usrID)

	// Squeeze our blurbs into the template, execute
	t.Execute(w, bs)
}
