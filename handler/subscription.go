package handler

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/adamsanghera/blurber/protobufs/dist/common"
	sub "github.com/adamsanghera/blurber/protobufs/dist/subscription"
)

func Subscribe(w http.ResponseWriter, r *http.Request) {
	log.Printf("HANDLERS-SUB: Request received")

	start := time.Now()
	defer func() {
		log.Printf("HANDLERS-SUB: Request serviced in %5.1f seconds", time.Since(start).Seconds())
	}()

	// Validate session
	validSesh, uname := validateSession(w, r)
	if !validSesh {
		return
	}

	// Only posts allowed
	log.Printf("%v", r.Method)
	if r.Method != "POST" {
		return
	}

	err := r.ParseForm()
	if err != nil {
		w.Write([]byte("Poorly-formed request"))
		return
	}

	// Obtain all variables
	leaderName := r.Form.Get("subscribe-leader")
	uid, err := userDB.GetUserID(uname)
	lid, lErr := userDB.GetUserID(leaderName)

	// Verify leader exists
	if lErr != nil {
		w.Write([]byte("The person you want to susbcribe to doesn't exist :(\n"))
		return
	}

	// Cannot self-subscribe
	if uid == lid {
		w.Write([]byte("You can't subscribe to yourself :(\n"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Submit subscription request to ledger
	_, err = subDB.Add(ctx, &sub.Subscription{
		Follower: &common.UserID{UserID: int32(uid)},
		Leader:   &common.UserID{UserID: int32(lid)},
	})

	if err != nil {
		panic(err)
	}

	_, err = blurbDB.InvalidateFeedCache(ctx, &common.UserID{UserID: int32(uid)})
	if err != nil {
		panic(err)
	}

	// Redirect to feed
	http.Redirect(w, r, "/feed/", http.StatusFound)
}

func Unsubscribe(w http.ResponseWriter, r *http.Request) {
	log.Printf("HANDLERS-UNSUB: Request received")

	start := time.Now()
	defer func() {
		log.Printf("HANDLERS-UNSUB: Request serviced in %5.1f seconds", time.Since(start).Seconds())
	}()

	// Validate session
	validSesh, uname := validateSession(w, r)
	if !validSesh {
		return
	}

	// Only posts allowed
	if r.Method != "POST" {
		log.Printf("whoops %v", r.Method)
		return
	}

	err := r.ParseForm()
	if err != nil {
		w.Write([]byte("Poorly-formed request"))
		return
	}

	// Obtain all variables
	leaderName := r.Form.Get("blurber")
	uid, _ := userDB.GetUserID(uname)
	lid, _ := userDB.GetUserID(leaderName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	subDB.Delete(ctx, &sub.Subscription{
		Follower: &common.UserID{UserID: int32(uid)},
		Leader:   &common.UserID{UserID: int32(lid)},
	})

	_, err = blurbDB.InvalidateFeedCache(ctx, &common.UserID{UserID: int32(uid)})
	if err != nil {
		panic(err)
	}

	// Redirect to feed
	http.Redirect(w, r, "/feed/", http.StatusFound)
}
