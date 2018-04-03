package handler

import (
	"log"
	"net/http"
)

func Subscribe(w http.ResponseWriter, r *http.Request) {
	log.Printf("SUB: URL is <%s>", r.URL.Path)

	validSesh, usernam := validateSession(w, r)
	if !validSesh {
		return
	}

	if r.Method != "POST" {
		return
	}

	err := r.ParseForm()
	if err != nil {
		w.Write([]byte("Poorly-formed request"))
		return
	}

	leaderField := r.Form.Get("subscribe-leader")

	uID, err := userDB.GetUserID(usernam)
	lID, lErr := userDB.GetUserID(leaderField)
	// Verify leader exists
	if lErr != nil {
		w.Write([]byte("The person you want to susbcribe doesn't exist :(\n"))
		return
	}
	// Check for subscription to self
	if uID == lID {
		w.Write([]byte("You can't subscribe to yourself :(\n"))
		return
	}

	subDB.AddSub(uID, lID)

	http.Redirect(w, r, "/feed/", http.StatusFound)
}
