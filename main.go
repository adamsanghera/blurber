package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"./blurb"
	reg "./registration"
)

var userDB = reg.NewLocalLedger()
var increasingCounter int

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("SERVER: URL is <%s>", r.URL.Path)
	if r.Method == "GET" {
		log.Printf("SERVER: Login – GET")
		http.FileServer(http.Dir("./static")).ServeHTTP(w, r)
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
			err := userDB.Add(reg.User{
				UID:  strconv.Itoa(increasingCounter),
				Name: r.Form.Get("reg-user"),
			}, r.Form.Get("reg-pass"))

			if err != nil {
				log.Printf("SERVER: Register – failed {%s}", err.Error())
				w.Write([]byte("Bad registration\n"))
			} else {
				w.Write([]byte("Good registration\n"))
			}
		} else {
			w.Write([]byte("Something went wrong\n"))
		}
	}
}

func FeedHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("FEED: URL is <%s>", r.URL.Path)
	tok, err := r.Cookie("token")
	if err != nil || tok.Value == "" {
		log.Printf("Redirecting to log in, tok {%s} ", tok)
		http.Redirect(w, r, "/login/", http.StatusFound)
		return
	}
	usr, err := r.Cookie("uname")
	if err != nil || usr.Value == "" {
		http.Redirect(w, r, "/login/", http.StatusFound)
		return
	}

	log.Printf("SERVER: Feed for %s with token %5s \n", usr.Value, tok.Value)
	err, token := userDB.Authorize(usr.Value, tok.Value)
	if err != nil {
		http.Redirect(w, r, "/login/", http.StatusFound)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(365 * 24 * time.Hour),
		Path:    "/",
	})

	t, err := template.ParseFiles("./static/feed/index.html", "./static/feed/index.css")
	if err != nil {
		panic(err)
	}

	bs := make([]blurb.Blurb, 3)

	bs[0] = blurb.Blurb{
		CreatorName: "Adam",
		Content:     "hi",
		Timestamp:   time.Now().Format("Jan 2 – 15:04 EDT"),
	}
	bs[1] = blurb.Blurb{
		CreatorName: "Adam",
		Content:     "bye",
		Timestamp:   time.Now().Format("Jan 2 – 15:04 EDT"),
	}
	bs[2] = blurb.Blurb{
		CreatorName: "Adam",
		Content:     "nooo",
		Timestamp:   time.Now().Format("Jan 2 – 15:04 EDT"),
	}

	t.Execute(w, bs)
	// http.FileServer(http.Dir("./static")).ServeHTTP(w, r)
}

func main() {
	log.Printf("SERVER: Starting...")
	increasingCounter = 0

	http.HandleFunc("/login/", LoginHandler)
	http.HandleFunc("/feed/", FeedHandler)

	http.ListenAndServe(":4000", nil)
}
