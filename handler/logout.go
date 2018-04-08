package handler

import (
    "log"
    "net/http"
    "time"
)

func Logout(w http.ResponseWriter, r *http.Request) {
    log.Printf("HANDLERS-LOGOUT: Request received")

    start := time.Now()
    defer func() {
        log.Printf("HANDLERS-LOGOUT: Request serviced in %5.1f seconds", time.Since(start).Seconds())
    }()

    // Check the cookie jar for a token
    cookieToken, err := r.Cookie("token")
    if err != nil || cookieToken.Value == "" {
        http.Redirect(w, r, "/login/", http.StatusFound)
        log.Printf("HANDLERS-LOGOUT: Failed & Redirected")
        return
    }

    // Check the cookie jar for a username
    cookieUsername, err := r.Cookie("uname")
    if err != nil || cookieUsername.Value == "" {
        http.Redirect(w, r, "/login/", http.StatusFound)
        log.Printf("HANDLERS-LOGOUT: Failed & Redirected")
        return
    }

    // Reset token
    err = userDB.CheckOut(cookieUsername.Value, cookieToken.Value)
    if err != nil {
        http.Redirect(w, r, "/login/expired/", http.StatusFound)
        log.Printf("HANDLERS-LOGOUT: Failed & Redirected")
        return
    }

    // Set cookies to expire
    cookieToken.Expires = time.Unix(0, 0)
    cookieUsername.Expires = time.Unix(0, 0)

    http.Redirect(w, r, "/login/", http.StatusFound)
}
