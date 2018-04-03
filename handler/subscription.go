package handler

import (
    "log"
    "net/http"
)

func Subscribe(w http.ResponseWriter, r *http.Request) {
    log.Printf("SUB: URL is <%s>", r.URL.Path)

    validSesh, usernam := validateSession(w, r)
    if !validSesh { return }

    if r.Method != "POST" { return }

    err := r.ParseForm() 
    if err != nil { panic(err) }

    leaderField := r.Form.Get("subscribe-leader")

    _, uID := userDB.GetUsrID(usernam)
    lErr, lID := userDB.GetUsrID(leaderField)
    // Verify leader exists
    if lErr != nil {
        w.Write([]byte("The person you want to susbcribe doesn't exist :(\n"))
        return
    }
    // Check for subscription to self
    if uID == lID {
        w.Write([]byte("You can't subscribe to yourself :(\n"))
        return
    }

    subDB.AddSub(usernam, leaderField)
    
    http.Redirect(w, r, "/feed/", http.StatusFound)
} 