package handler

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/adamsanghera/blurber/protobufs/dist/blurb"
	"github.com/adamsanghera/blurber/protobufs/dist/user"
)

func Feed(w http.ResponseWriter, r *http.Request) {
	log.Printf("HANDLERS-FEED: Request received")

	start := time.Now()
	defer func() {
		log.Printf("HANDLERS-FEED: Request serviced in %5.1f seconds", time.Since(start).Seconds())
	}()

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

	// Get our UID
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uid, err := userDB.GetID(ctx, &user.Username{Username: uname})
	if err != nil {
		w.Write([]byte("Something went very wrong"))
		panic(err)
	}

	// Get the leader map for UID
	leaderSet, err := subDB.GetLeadersOf(ctx, uid)
	if err != nil {
		w.Write([]byte("Something went very wrong"))
		panic(err)
	}

	// ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	// Generate the feed
	bs, err := blurbDB.GenerateFeed(ctx, &blurb.FeedParameters{
		RequestorID: uid,
		LeaderIDs:   leaderSet.Leaders,
	})

	if err != nil {
		panic(err)
	}

	// Squeeze our blurbs into the template, execute
	t.Execute(w, bs.Blurbs)
}
