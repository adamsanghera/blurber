package handler

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/adamsanghera/blurber/protobufs/dist/blurb"
	"github.com/adamsanghera/blurber/protobufs/dist/common"
)

func Blurb(w http.ResponseWriter, r *http.Request) {
	log.Printf("HANDLERS-BLURB: %s request received", r.Method)

	start := time.Now()
	defer func() {
		log.Printf("HANDLERS-BLURB: Request serviced in %5.1f seconds", time.Since(start).Seconds())
	}()

	validSesh, username := validateSession(w, r)

	if !validSesh {
		return
	}

	if r.Method == "GET" {
		// Build template
		t, err := template.ParseFiles("./static-assets/blurb/index.html")
		if err != nil {
			panic(err)
		}

		t.Execute(w, nil)
		return
	}

	if r.Method == "POST" {

		// Parse the form
		err := r.ParseForm()
		if err != nil {
			panic(err)
		}

		if r.URL.Path == "/blurb/add" {
			// Retrieve uid
			usrID, err := userDB.GetUserID(username)
			if err != nil {
				w.Write([]byte("Something went wrong\n\terr: " + err.Error()))
				return
			}

			// Create blurb
			// TODO: Validate content to be non-empty
			content := r.Form.Get("blurb-write-text")

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			_, err = blurbDB.Add(ctx, &blurb.NewBlurb{
				Author:   &common.UserID{UserID: int32(usrID)},
				Content:  content,
				Username: username,
			})

			if err != nil {
				panic(err)
			}

			http.Redirect(w, r, "/profile/", http.StatusFound)
			return
		}

		if r.URL.Path == "/blurb/remove" {
			// Retrieve uid
			usrID, err := userDB.GetUserID(username)
			if err != nil {
				w.Write([]byte("Something went wrong\n\terr: " + err.Error()))
				return
			}

			sBid := r.Form.Get("remove-bid")
			bid, _ := strconv.Atoi(sBid)

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			blurbDB.Delete(ctx, &blurb.BlurbIndex{
				Author:  &common.UserID{UserID: int32(usrID)},
				BlurbID: int32(bid),
			})

			http.Redirect(w, r, "/profile/", http.StatusFound)
			return
		}

		w.Write([]byte("Something went wrong\n"))
	}
}
