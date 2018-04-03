package handler

import (
    "log"
    "net/http"
)

func RemoveAcc(w http.ResponseWriter, r *http.Request) {
    log.Printf("LOGIN HANDLER: URL is <%s>", r.URL.Path)    
    
    validSesh, username := validateSession(w, r)
    if !validSesh { return }

    _, usrID := userDB.GetUsrID(username)

    userDB.Remove(username)
    subDB.RemoveUsr(usrID)
    lbl.RemoveAllBlurb(usrID)

    http.Redirect(w, r, "/login/", http.StatusFound)
}
