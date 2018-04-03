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

var lbl = blurb.NewLocalLedger()
var bidCounter int

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
	err, token := userDB.CheckIn(cookieUsername.Value, cookieToken.Value)
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
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
			err := userDB.Add(reg.User{
				UID:  strconv.Itoa(increasingCounter),
				Name: r.Form.Get("reg-user"),
			}, r.Form.Get("reg-pass"))

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

func FeedHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("FEED: URL is <%s>", r.URL.Path)

	// Validate and renew our cookies
	validSesh, uname := validateSession(w, r)
	if !validSesh {
		return
	}

	// Build template
	t, err := template.ParseFiles("./static-assets/feed/index.html")
	if err != nil {
		panic(err)
	}

	// TODO: Implement real logic here.
	// Obtain our blurb list
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

	err, usrID := userDB.GetUsrID(uname)
	bs = append(bs, lbl.GetUsrBlurb(usrID)...)

	// Squeeze our blurbs into the template, execute
	t.Execute(w, bs)
}

func BlurbHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("BLURB: URL is <%s>", r.URL.Path)
	if r.Method == "GET" {
		log.Printf("BLURB: GET")
		http.FileServer(http.Dir("./static-assets")).ServeHTTP(w, r)
	} else {
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}
		if r.URL.Path == "/blurb/add" {
			// Get uid
			usr, err := r.Cookie("uname")
			if err != nil || usr.Value == "" {
				http.Redirect(w, r, "/login/", http.StatusFound)
				return
			}

			err, usrID := userDB.GetUsrID(usr.Value)

			// Add blurb
			content := r.Form.Get("burb-text")
			newBlurb := blurb.Blurb{
				Content:     content,
				Timestamp:   time.Now().Format("Jan 2 – 15:04 EDT"),
				BID:         strconv.Itoa(bidCounter),
				CreatorName: usr.Value,
			}
			bidCounter += 1

			lbl.AddBlurb(usrID, newBlurb)

			log.Printf("BLURB: User %v (id %v) - New blurb added: %v", usr.Value, usrID, content)

			http.Redirect(w, r, "/feed/", http.StatusFound)
		} else {
			w.Write([]byte("Something went wrong\n"))
		}
	}
}

func main() {
	log.Printf("SERVER: Starting...")
	increasingCounter = 0

	// Dev only - for quick loging
	userDB.Add(reg.User{"dev", "-1"}, "root")

	http.HandleFunc("/login/", LoginHandler)
	http.HandleFunc("/feed/", FeedHandler)
	http.HandleFunc("/blurb/", BlurbHandler)

	http.ListenAndServe(":4000", nil)
}