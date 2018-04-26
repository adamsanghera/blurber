package handler

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/adamsanghera/blurber-protobufs/dist/user"
)

func RemoveAcc(w http.ResponseWriter, r *http.Request) {
	log.Printf("HANDLERS-REMOVE-ACC: Request received")

	start := time.Now()
	defer func() {
		log.Printf("HANDLERS-REMOVE-ACC: Request serviced in %5.1f seconds", time.Since(start).Seconds())
	}()

	// Validate session
	validSesh, username := validateSession(w, r)
	if !validSesh {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Retrieve uid
	usrID, err := userDB.GetID(ctx, &user.Username{Username: username})
	if err != nil {
		w.Write([]byte("Something went very wrong\n\tError Message: " + err.Error()))
		return
	}

	_, err = userDB.Delete(ctx, &user.Username{Username: username})
	if err != nil {
		panic(err)
	}

	_, err = subDB.DeletePresenceOf(ctx, usrID)
	if err != nil {
		panic(err)
	}

	_, err = blurbDB.DeleteHistoryOf(ctx, usrID)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/login/", http.StatusFound)
}
