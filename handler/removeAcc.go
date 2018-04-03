package handler

import (
	"log"
	"net/http"
	"time"
)

func RemoveAcc(w http.ResponseWriter, r *http.Request) {
	log.Printf("HANDLERS-REMOVE-ACC: Request received")

	start := time.Now()
	defer func() {
		log.Printf("HANDLERS-REMOVE-ACC: Request serviced in %5.1f seconds", time.Since(start).Seconds())
	}()

	// Validate session
	validSesh, username := validateSession(w, r)
	if !validSesh {
		return
	}

	// Retrieve UID
	usrID, err := userDB.GetUserID(username)
	if err != nil {
		w.Write([]byte("Something went very wrong\n\tError Message: " + err.Error()))
		return
	}

	userDB.Remove(username)
	subDB.RemoveUser(usrID)
	blurbDB.RemoveAllBlurbsBy(usrID)

	http.Redirect(w, r, "/login/", http.StatusFound)
}
