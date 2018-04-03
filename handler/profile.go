package handler

import (
	"html/template"
	"net/http"

	"../blurb"
)

// Profile is the handler for /profile/ requests
func Profile(w http.ResponseWriter, req *http.Request) {
	isValid, username := validateSession(w, req)
	if !isValid {
		return
	}

	// Build templates
	t, err := template.ParseFiles("./static-assets/feed/profile.html")
	if err != nil {
		panic(err)
	}

	// Get the UID, so that we can retrieve this user's blurb
	err, uid := userDB.GetUsrID(username)
	if err != nil {
		w.Write([]byte("Something went wrong while retrieving user blurbs"))
	}

	data := struct {
		name     string
		bio      string
		username string
		blurbs   []blurb.Blurb
	}{
		username,
		"<No name yet>",
		"<No bio yet>",
		lbl.GetUsrBlurb(uid),
	}

	t.Execute(w, data)
}
