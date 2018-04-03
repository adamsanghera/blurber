package handler

import (
	"log"
	"net/http"
	"time"
)

func Subscribe(w http.ResponseWriter, r *http.Request) {
	log.Printf("HANDLERS-SUB: Request received")

	start := time.Now()
	defer func() {
		log.Printf("HANDLERS-SUB: Request serviced in %5.1f seconds", time.Since(start).Seconds())
	}()

	// Validate session
	validSesh, uname := validateSession(w, r)
	if !validSesh {
		return
	}

	// No posts allowed
	if r.Method != "POST" {
		return
	}

	err := r.ParseForm()
	if err != nil {
		w.Write([]byte("Poorly-formed request"))
		return
	}

	// Obtain all variables
	leaderName := r.Form.Get("subscribe-leader")
	uID, err := userDB.GetUserID(uname)
	lID, lErr := userDB.GetUserID(leaderName)

	// Verify leader exists
	if lErr != nil {
		w.Write([]byte("The person you want to susbcribe to doesn't exist :(\n"))
		return
	}

	// Cannot self-subscribe
	if uID == lID {
		w.Write([]byte("You can't subscribe to yourself :(\n"))
		return
	}

	// Submit subscription request to ledger
	subDB.AddSub(uID, lID)

	// Redirect to feed
	http.Redirect(w, r, "/feed/", http.StatusFound)
}
