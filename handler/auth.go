package handler

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/adamsanghera/blurber-protobufs/dist/user"
)

func Auth(w http.ResponseWriter, r *http.Request) {
	log.Printf("HANDLERS-AUTH: %s-request received", r.Method)

	start := time.Now()
	defer func() {
		log.Printf("HANDLERS-AUTH: Request serviced in %5.1f seconds", time.Since(start).Seconds())
	}()

	errMsg := struct {
		ErrMsg string
	}{}

	t, err := template.ParseFiles("./static-assets/login/index.html")
	if err != nil {
		panic(err)
	}

	if r.Method == "GET" {
		if r.URL.Path == "/login/expired/" {
			errMsg.ErrMsg = "Sorry, your session expired"
			t.Execute(w, errMsg)
			return
		}

		t.Execute(w, nil)
		return
	}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}

		var username string
		var password string

		// We're doing something a bit clever here.
		// Extract username/password depending on request type, and returning "early" with errors if encountered
		if r.URL.Path == "/login/existing" {
			// Login
			log.Printf("HANDLERS-AUTH: Log-in request received")
			username = r.Form.Get("login-user")
			password = r.Form.Get("login-pass")

		} else if r.URL.Path == "/login/new" {
			// Registration
			log.Printf("HANDLERS-AUTH: Registration request received")
			username = r.Form.Get("reg-user")
			password = r.Form.Get("reg-pass")

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			_, err := userDB.Add(ctx, &user.Credentials{Username: username, Password: password})

			if err != nil {
				errMsg.ErrMsg = err.Error()
				t.Execute(w, errMsg)
				return
			}
		} else {
			errMsg.ErrMsg = "Invalid request"
			t.Execute(w, errMsg)
			return
		}

		// Actually log the user in
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		token, err := userDB.LogIn(ctx, &user.Credentials{Username: username, Password: password})
		if err != nil {
			errMsg.ErrMsg = err.Error()
			t.Execute(w, errMsg)
			return
		}

		// New cookies for token and username
		tokCook := &http.Cookie{
			Name:    "token",
			Value:   token.Token,
			Expires: time.Now().Add(24 * time.Hour),
			Path:    "/",
		}

		nameCook := &http.Cookie{
			Name:    "uname",
			Value:   username,
			Expires: time.Now().Add(24 * time.Hour),
			Path:    "/",
		}

		// Set the cookies
		http.SetCookie(w, tokCook)
		http.SetCookie(w, nameCook)

		http.Redirect(w, r, "/feed", http.StatusFound)
		return
	}
	errMsg.ErrMsg = "Something went wrong"
	t.Execute(w, errMsg)
}
