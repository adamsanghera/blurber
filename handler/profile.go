package handler

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/adamsanghera/blurber/protobufs/dist/blurb"
	"github.com/adamsanghera/blurber/protobufs/dist/user"
)

// Profile is the handler for /profile/ requests
func Profile(w http.ResponseWriter, req *http.Request) {
	log.Printf("HANDLERS-PROFILE: Request received")

	start := time.Now()
	defer func() {
		log.Printf("HANDLERS-PROFILE: Request serviced in %5.1f seconds", time.Since(start).Seconds())
	}()

	// Validate the sesh
	isValid, username := validateSession(w, req)
	if !isValid {
		return
	}

	// Build templates
	t, err := template.ParseFiles("./static-assets/profile/index.html")
	if err != nil {
		panic(err)
	}

	// Get the UID, so that we can retrieve this user's blurb
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Retrieve uid
	uid, err := userDB.GetID(ctx, &user.Username{Username: username})
	if err != nil {
		w.Write([]byte("Something went very wrong"))
	}

	// Obtain and sort blurb
	blurbs, err := blurbDB.GetRecentBy(ctx, uid)
	if err != nil {
		panic(err)
	}

	// Create a data packet to send back
	data := struct {
		Name     string
		Bio      string
		Username string
		Blurbs   []*blurb.Blurb
	}{
		"<No name yet>",
		"<No bio yet>",
		username,
		blurbs.Blurbs,
	}

	t.Execute(w, data)
}
