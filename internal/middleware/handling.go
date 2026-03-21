package middleware

import (
	"log"
	"net/http"
)

func HandleFileServerWithLogging(w http.ResponseWriter, r *http.Request) {
	log.Println("serving static content at " + r.URL.Path)
	http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))).ServeHTTP(w, r)
}
