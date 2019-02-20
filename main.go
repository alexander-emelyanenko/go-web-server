package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type notFoundHandler struct{}

func (n *notFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("<h1>404: The page you're looking for doesn't exists</h1>"))
}

func home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<h1>Welcome to main page of our web site!</h1>"))
}

func contact(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<h1>Soon here will be our contacts... I hope</h1>"))
}

func faq(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<h1>Frequently asked questions</h1>"))
}

func main() {
	router := httprouter.New()
	router.GET("/", home)
	router.GET("/contact", contact)
	router.GET("/faq", faq)
	router.NotFound = &notFoundHandler{}
	log.Fatal(http.ListenAndServe(":3000", router))
}
