package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/alexander-emelyanenko/go-web-server/controllers"
	"github.com/gorilla/mux"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "qwerty"
	dbname   = "go-web-server-dev"
)

func main() {
	// psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	staticController := controllers.NewStatic()
	usersController := controllers.NewUsers()

	router := mux.NewRouter()
	router.Handle("/", staticController.Home).Methods("GET")
	router.Handle("/contact", staticController.Contact).Methods("GET")
	router.HandleFunc("/signup", usersController.New).Methods("GET")
	router.HandleFunc("/signup", usersController.Create).Methods("POST")

	log.Fatal(http.ListenAndServe(":3000", router))
}