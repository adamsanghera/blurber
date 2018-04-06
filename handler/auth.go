package handler

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

func Auth(w http.ResponseWriter, r *http.Request) {
	log.Printf("HANDLERS-AUTH: %s-request received", r.Method)

	start := time.Now()
	defer func() {
		log.Printf("HANDLERS-AUTH: Request serviced in %5.1f seconds", time.Since(start).Seconds())
	}()

	if r.Method == "GET" {
		t, err := template.ParseFiles("./static-assets/login/index.html")
		if err != nil {
			panic(err)
		}

		msg := struct {
			ErrMsg string
		}{}
		if r.URL.Path == "/login/expired/" {
			msg.ErrMsg = "Sorry, your session expired"
			t.Execute(w, msg)
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

			err := userDB.AddNewUser(username, password)

			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}
		} else {
			w.Write([]byte("Something went wrong\n"))
			return
		}

		// Actually log the user in

		token, err := userDB.LogIn(username, password)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		// New cookies for token and username
		tokCook := &http.Cookie{
			Name:    "token",
			Value:   token,
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

	w.Write([]byte("Something went wrong\n"))
}
