package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alexander-emelyanenko/go-web-server/controllers"
	"github.com/alexander-emelyanenko/go-web-server/models"
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
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// Initializing our UserService
	services, err := models.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.AutoMigrate()

	staticController := controllers.NewStatic()
	usersController := controllers.NewUsers(services.User)
	galleriesController := controllers.NewGalleries(services.Gallery)

	router := mux.NewRouter()

	router.Handle("/", staticController.Home).Methods("GET")
	router.Handle("/contact", staticController.Contact).Methods("GET")

	router.HandleFunc("/signup", usersController.New).Methods("GET")
	router.HandleFunc("/signup", usersController.Create).Methods("POST")

	router.Handle("/login", usersController.LoginView).Methods("GET")
	router.HandleFunc("/login", usersController.Login).Methods("POST")

	router.Handle("/galleries/new", galleriesController.New).Methods("GET")

	router.HandleFunc("/cookietest", usersController.CookieTest).Methods("GET")

	log.Fatal(http.ListenAndServe(":3000", router))
}
