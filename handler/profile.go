package handler

import (
	"html/template"
	"log"
	"net/http"

	"../blurb"
)

// Profile is the handler for /profile/ requests
func Profile(w http.ResponseWriter, req *http.Request) {
	log.Printf("PROFILE: Validating request")
	isValid, username := validateSession(w, req)
	if !isValid {
		return
	}

	log.Printf("PROFILE: Request validated for %s", username)
	// Build templates
	t, err := template.ParseFiles("./static-assets/profile/index.html")
	if err != nil {
		panic(err)
	}

	// Get the UID, so that we can retrieve this user's blurb
	err, uid := userDB.GetUsrID(username)
	if err != nil {
		w.Write([]byte("Something went wrong while retrieving user blurbs"))
	}

	data := struct {
		Name     string
		Bio      string
		Username string
		blurbs   []blurb.Blurb
	}{
		"<No name yet>",
		"<No bio yet>",
		username,
		lbl.GetUsrBlurb(uid),
	}

	t.Execute(w, data)
}
