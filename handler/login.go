package handler

import (
	"log"
	"net/http"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {
	log.Printf("SERVER: URL is <%s>", r.URL.Path)
	if r.Method == "GET" {
		log.Printf("SERVER: Login – GET")
		http.FileServer(http.Dir("./static-assets")).ServeHTTP(w, r)
	} else {
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}
		if r.URL.Path == "/login/existing" {
			log.Printf("SERVER: Login – Attempted login")
			err, token := userDB.LogIn(r.Form.Get("login-user"), r.Form.Get("login-pass"))
			if err != nil {
				log.Printf("SERVER: Login – login error {%s}", err.Error())
				w.Write([]byte("Bad password\n"))
			} else {
				tokCook := &http.Cookie{
					Name:    "token",
					Value:   token,
					Expires: time.Now().Add(365 * 24 * time.Hour),
					Path:    "/",
				}
				nameCook := &http.Cookie{
					Name:    "uname",
					Value:   r.Form.Get("login-user"),
					Expires: time.Now().Add(365 * 24 * time.Hour),
					Path:    "/",
				}
				http.SetCookie(w, tokCook)
				http.SetCookie(w, nameCook)
				log.Printf("SERVER: Login – Good login, redirecting to feed")
				http.Redirect(w, r, "/feed/", http.StatusFound)
			}
		} else if r.URL.Path == "/login/new" {
			log.Printf("SERVER: Register – New attempt with {%s, %s}", r.Form["reg-user"], r.Form["reg-pass"])
			err := userDB.AddNewUser(r.Form.Get("reg-user"), r.Form.Get("reg-pass"))

			if err != nil {
				log.Printf("SERVER: Register – failed {%s}", err.Error())
				w.Write([]byte("Bad registration\n"))
			} else {
				w.Write([]byte("Good registration\n"))

				// Should redirect new user to their feed too?
				// Also increment increasingCounter

			}
		} else {
			w.Write([]byte("Something went wrong\n"))
		}
	}
}
