package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"../blurb"
)

func Blurb(w http.ResponseWriter, r *http.Request) {
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
