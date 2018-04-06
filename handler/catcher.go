package handler

import "net/http"

func Catcher(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "/login/", http.StatusFound)
}
