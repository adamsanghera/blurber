package handler

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/adamsanghera/blurber/protobufs/dist/common"
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

	// Retrieve UID
	usrID, err := userDB.GetUserID(username)
	if err != nil {
		w.Write([]byte("Something went very wrong\n\tError Message: " + err.Error()))
		return
	}

	userDB.Remove(username)
	subDB.RemoveUser(usrID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = blurbDB.DeleteHistoryOf(ctx, &common.UserID{UserID: int32(usrID)})
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/login/", http.StatusFound)
}
