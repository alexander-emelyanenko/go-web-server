package main

import (
	"log"
	"net/http"

	"github.com/alexander-emelyanenko/go-web-server/controllers"
	"github.com/alexander-emelyanenko/go-web-server/views"
	"github.com/gorilla/mux"
)

var homeView *views.View
var contactView *views.View

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(homeView.Render(w, nil))
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(contactView.Render(w, nil))
}

func main() {
	homeView = views.NewView("bootstrap", "views/home.gohtml")
	contactView = views.NewView("bootstrap", "views/contact.gohtml")

	usersController := controllers.NewUsers()

	router := mux.NewRouter()
	router.HandleFunc("/", home)
	router.HandleFunc("/contact", contact)
	router.HandleFunc("/signup", usersController.New)

	log.Fatal(http.ListenAndServe(":3000", router))
}
