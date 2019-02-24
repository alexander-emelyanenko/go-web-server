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

	// Connecting to DB
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	us.DestructiveReset()

	// Creating new user
	user := &models.User{
		Name:  "Alexander Ivanov",
		Email: "ivanov@gmail.com",
	}

	if err := us.Create(user); err != nil {
		panic(err)
	}

	// Fetching user from DB by Email
	newUser, err := us.ByEmail("ivanov@gmail.com")
	if err != nil {
		panic(err)
	}

	fmt.Printf("New user is %v\n", *newUser)

	staticController := controllers.NewStatic()
	usersController := controllers.NewUsers()

	router := mux.NewRouter()
	router.Handle("/", staticController.Home).Methods("GET")
	router.Handle("/contact", staticController.Contact).Methods("GET")
	router.HandleFunc("/signup", usersController.New).Methods("GET")
	router.HandleFunc("/signup", usersController.Create).Methods("POST")

	log.Fatal(http.ListenAndServe(":3000", router))
}
