package main

import (
	"log"
	"net/http"

	"github.com/alexander-emelyanenko/go-web-server/views"

	"github.com/julienschmidt/httprouter"
)

var homeView *views.View
var contactView *views.View

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	must(homeView.Render(w, nil))
}

func contact(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	must(contactView.Render(w, nil))
}

func main() {
	homeView = views.NewView("bootstrap", "views/home.gohtml")
	contactView = views.NewView("bootstrap", "views/contact.gohtml")

	router := httprouter.New()
	router.GET("/", home)
	router.GET("/contact", contact)

	log.Fatal(http.ListenAndServe(":3000", router))
}
