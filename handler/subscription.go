package handler

import (
	"context"
	"log"
	"net/http"
	"time"

	sub "github.com/adamsanghera/blurber/protobufs/dist/subscription"
	"github.com/adamsanghera/blurber/protobufs/dist/user"
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Retrieve uid
	uid, err := userDB.GetID(ctx, &user.Username{Username: uname})
	lid, lErr := userDB.GetID(ctx, &user.Username{Username: leaderName})

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

	// Submit subscription request to ledger
	_, err = subDB.Add(ctx, &sub.Subscription{
		Follower: uid,
		Leader:   lid,
	})

	if err != nil {
		panic(err)
	}

	_, err = blurbDB.InvalidateFeedCache(ctx, uid)
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Retrieve uid
	uid, _ := userDB.GetID(ctx, &user.Username{Username: uname})
	lid, _ := userDB.GetID(ctx, &user.Username{Username: leaderName})

	subDB.Delete(ctx, &sub.Subscription{
		Follower: uid,
		Leader:   lid,
	})

	_, err = blurbDB.InvalidateFeedCache(ctx, uid)
	if err != nil {
		panic(err)
	}

	// Redirect to feed
	http.Redirect(w, r, "/feed/", http.StatusFound)
}
