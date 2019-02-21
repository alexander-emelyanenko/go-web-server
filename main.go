package main

import (
	"log"
	"net/http"

	"github.com/alexander-emelyanenko/go-web-server/controllers"
	"github.com/gorilla/mux"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	staticController := controllers.NewStatic()
	usersController := controllers.NewUsers()

	router := mux.NewRouter()
	router.Handle("/", staticController.Home).Methods("GET")
	router.Handle("/contact", staticController.Contact).Methods("GET")
	router.HandleFunc("/signup", usersController.New).Methods("GET")
	router.HandleFunc("/signup", usersController.Create).Methods("POST")

	log.Fatal(http.ListenAndServe(":3000", router))
}
