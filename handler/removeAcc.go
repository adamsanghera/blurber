package handler

import (
	"log"
	"net/http"
)

func RemoveAcc(w http.ResponseWriter, r *http.Request) {
	log.Printf("LOGIN HANDLER: URL is <%s>", r.URL.Path)

	validSesh, username := validateSession(w, r)
	if !validSesh {
		return
	}

	usrID, err := userDB.GetUserID(username)
	if err != nil {
		panic(err)
	}

	userDB.Remove(username)
	subDB.RemoveUser(usrID)
	lbl.RemoveAllBlurbsBy(usrID)

	http.Redirect(w, r, "/login/", http.StatusFound)
}
