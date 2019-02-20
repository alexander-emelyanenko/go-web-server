package main

import (
	"log"
	"net/http"

	"github.com/alexander-emelyanenko/go-web-server/views"

	"github.com/julienschmidt/httprouter"
)

var homeView *views.View
var contactView *views.View

func home(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	if err := homeView.Template.Execute(w, nil); err != nil {
		panic(err)
	}
}

func contact(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	if err := contactView.Template.Execute(w, nil); err != nil {
		panic(err)
	}
}

func main() {
	homeView = views.NewView("views/home.gohtml")
	contactView = views.NewView("views/contact.gohtml")

	router := httprouter.New()
	router.GET("/", home)
	router.GET("/contact", contact)

	log.Fatal(http.ListenAndServe(":3000", router))
}
