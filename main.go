package main

import (
	"net/http"
)

func handleFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if r.URL.Path == "/" {
		w.Write([]byte("<h1>Welcome to main page of our web site!</h1>"))
	} else if r.URL.Path == "/contact" {
		w.Write([]byte("<h1>Soon here will be our contacts... I hope</h1>"))
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("<h1>404: The page you are looking doesn't exists :(</h1>"))
	}
}

func main() {
	http.HandleFunc("/", handleFunc)
	http.ListenAndServe(":3000", nil)
}
