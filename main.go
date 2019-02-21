package main

import (
	"log"
	"net/http"

	"github.com/alexander-emelyanenko/go-web-server/views"

	"github.com/julienschmidt/httprouter"
)

var homeView *views.View
var contactView *views.View
var signupView *views.View

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

func signup(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	must(signupView.Render(w, nil))
}

func main() {
	homeView = views.NewView("bootstrap", "views/home.gohtml")
	contactView = views.NewView("bootstrap", "views/contact.gohtml")
	signupView = views.NewView("bootstrap", "views/signup.gohtml")

	router := httprouter.New()
	router.GET("/", home)
	router.GET("/contact", contact)
	router.GET("/signup", signup)

	log.Fatal(http.ListenAndServe(":3000", router))
}
