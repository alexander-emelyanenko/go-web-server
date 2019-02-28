package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alexander-emelyanenko/go-web-server/controllers"
	"github.com/alexander-emelyanenko/go-web-server/middleware"
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

	router := mux.NewRouter()

	staticController := controllers.NewStatic()
	usersController := controllers.NewUsers(services.User)
	galleriesController := controllers.NewGalleries(services.Gallery, router)

	requireUserMw := middleware.RequireUser{
		UserService: services.User,
	}

	router.Handle("/", staticController.Home).Methods("GET")
	router.Handle("/contact", staticController.Contact).Methods("GET")

	router.HandleFunc("/signup", usersController.New).Methods("GET")
	router.HandleFunc("/signup", usersController.Create).Methods("POST")

	router.Handle("/login", usersController.LoginView).Methods("GET")
	router.HandleFunc("/login", usersController.Login).Methods("POST")

	router.Handle("/galleries/new", requireUserMw.Apply(galleriesController.New)).Methods("GET")
	router.HandleFunc("/galleries", requireUserMw.ApplyFn(galleriesController.Create)).Methods("POST")
	router.HandleFunc("/galleries/{id:[0-9]+}", galleriesController.Show).Methods("GET").Name(controllers.ShowGallery)
	router.HandleFunc("/galleries/{id:[0-9]+}/edit", requireUserMw.ApplyFn(galleriesController.Edit)).Methods("GET")
	router.HandleFunc("/galleries/{id:[0-9]+}/update", requireUserMw.ApplyFn(galleriesController.Update)).Methods("POST")
	router.HandleFunc("/galleries/{id:[0-9]+}/delete", requireUserMw.ApplyFn(galleriesController.Delete)).Methods("POST")

	router.HandleFunc("/cookietest", usersController.CookieTest).Methods("GET")

	log.Fatal(http.ListenAndServe(":3000", router))
}
