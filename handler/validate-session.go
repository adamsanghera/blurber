package handler

import (
	"log"
	"net/http"
	"time"
)

func validateSession(w http.ResponseWriter, r *http.Request) (bool, string) {
	// Check the cookie jar for a token
	cookieToken, err := r.Cookie("token")
	if err != nil || cookieToken.Value == "" {
		log.Printf("Redirecting to log in, tok {%s} ", cookieToken.Value)
		http.Redirect(w, r, "/login/", http.StatusFound)
		return false, ""
	}

	// Check the cookie jar for a username
	cookieUsername, err := r.Cookie("uname")
	if err != nil || cookieUsername.Value == "" {
		http.Redirect(w, r, "/login/", http.StatusFound)
		return false, ""
	}

	// Make sure that the token matches the username
	token, err := userDB.CheckIn(cookieUsername.Value, cookieToken.Value)
	if err != nil {
		http.Redirect(w, r, "/login/", http.StatusFound)
		return false, ""
	}

	// Set the token to be the new value
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(365 * 24 * time.Hour),
		Path:    "/",
	})

	return true, cookieUsername.Value
}
