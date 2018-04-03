package handler

import (
	"html/template"
	"log"
	"net/http"
)

func Blurb(w http.ResponseWriter, r *http.Request) {
	log.Printf("BLURB: URL is <%s>", r.URL.Path)

	validSesh, username := validateSession(w, r)
	if !validSesh {
		return
	}

	if r.Method == "GET" {
		log.Printf("BLURB: GET")

		// Build template
		t, err := template.ParseFiles("./static-assets/blurb/index.html")
		if err != nil {
			panic(err)
		}

		t.Execute(w, nil)

	} else {
		// Parse the form
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}

		if r.URL.Path == "/blurb/add" {
			// Determine uid
			err, usrID := userDB.GetUsrID(username)
			if err != nil {
				w.Write([]byte("Something went wrong\n"))
				return
			}

			// Create blurb
			// TODO: Validate content to be non-empty
			content := r.Form.Get("burb-text")
			lbl.AddNewBlurb(usrID, content, username)

			log.Printf("BLURB: User %v (id %v) - New blurb added: %v", username, usrID, content)

			http.Redirect(w, r, "/feed/", http.StatusFound)
		} else {
			w.Write([]byte("Something went wrong\n"))
		}
	}
}
