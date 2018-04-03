package handler

import (
    "log"
    "net/http"
)

func RemoveAcc(w http.ResponseWriter, r *http.Request) {
    log.Printf("LOGIN HANDLER: URL is <%s>", r.URL.Path)    
    
    validSesh, username := validateSession(w, r)
    if !validSesh { return }

    if r.Method != "DELETE" { return }

    if r.URL.Path != "/removeAcc" { return }

    _, usrID := userDB.GetUsrID(username)

    userDB.Remove(username)
    subDB.RemoveUsr(usrID)
    lbl.RemoveAllBlurb(usrID)
}
