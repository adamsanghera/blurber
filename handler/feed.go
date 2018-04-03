package handler

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"../blurb"
)

func Feed(w http.ResponseWriter, r *http.Request) {
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
