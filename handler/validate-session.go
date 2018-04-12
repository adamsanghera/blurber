package handler

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/adamsanghera/blurber/protobufs/dist/user"
)

func validateSession(w http.ResponseWriter, r *http.Request) (bool, string) {
	log.Printf("HANDLERS-VALIDATE: Validating session")

	start := time.Now()
	defer func() {
		log.Printf("HANDLERS-VALIDATE: Request serviced in %5.1f seconds", time.Since(start).Seconds())
	}()

	// Check the cookie jar for a token
	cookieToken, err := r.Cookie("token")
	if err != nil || cookieToken.Value == "" {
		http.Redirect(w, r, "/login/", http.StatusFound)
		log.Printf("HANDLERS-VALIDATE: Failed & Redirected")
		return false, ""
	}

	// Check the cookie jar for a username
	cookieUsername, err := r.Cookie("uname")
	if err != nil || cookieUsername.Value == "" {
		http.Redirect(w, r, "/login/", http.StatusFound)
		log.Printf("HANDLERS-VALIDATE: Failed & Redirected")
		return false, ""
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Make sure that the token matches the username
	token, err := userDB.CheckIn(ctx, &user.SessionCredentials{
		Username: cookieUsername.Value,
		Token:    cookieToken.Value})
	if err != nil {
		http.Redirect(w, r, "/login/expired/", http.StatusFound)
		log.Printf("HANDLERS-VALIDATE: Failed & Redirected")
		return false, ""
	}

	// Set the token to be the new value
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token.Token,
		Expires: time.Now().Add(365 * 24 * time.Hour),
		Path:    "/",
	})

	log.Printf("HANDLERS-VALIDATE: Passed")

	return true, cookieUsername.Value
}
