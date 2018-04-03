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
		t.Execute(w, nil)
		return
	}

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}

		if r.URL.Path == "/login/existing" {
			// Login
			log.Printf("HANDLERS-AUTH: Log-in request received")

		} else if r.URL.Path == "/login/new" {
			// Registration
			log.Printf("HANDLERS-AUTH: Registration request received")
			err := userDB.AddNewUser(r.Form.Get("reg-user"), r.Form.Get("reg-pass"))

			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}
		} else {
			w.Write([]byte("Something went wrong\n"))
			return
		}

		// Actually log the user in

		token, err := userDB.LogIn(r.Form.Get("login-user"), r.Form.Get("login-pass"))
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
			Value:   r.Form.Get("login-user"),
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
