package handler

import (
	"html/template"
	"log"
	"net/http"
	"sort"

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
	uid, err := userDB.GetUserID(username)
	if err != nil {
		w.Write([]byte("Something went wrong while retrieving user blurbs"))
	}

	blurbs := lbl.GetBlurbsCreatedBy(uid)

	sort.Slice(blurbs, func(i, j int) bool {
		return blurbs[i].Time.Before(blurbs[j].Time)
	})

	data := struct {
		Name     string
		Bio      string
		Username string
		Blurbs   []blurb.Blurb
	}{
		"<No name yet>",
		"<No bio yet>",
		username,
		blurbs,
	}

	log.Printf("PROFILE: User has written %d blurbs", len(blurbs))

	t.Execute(w, data)
}
