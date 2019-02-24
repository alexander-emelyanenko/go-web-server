package controllers

import (
	"fmt"
	"net/http"

	"github.com/alexander-emelyanenko/go-web-server/models"
	"github.com/alexander-emelyanenko/go-web-server/views"
)

// Users struct describes our users controller
type Users struct {
	NewView     *views.View
	userService *models.UserService
}

// NewUsers method returns Users struct
func NewUsers(userService *models.UserService) *Users {
	return &Users{
		NewView:     views.NewView("bootstrap", "users/new"),
		userService: userService,
	}
}

// SignupForm describes sign up request
type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// New method handles sign up request
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

// Create new user
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}

	user := models.User{
		Name:  form.Name,
		Email: form.Email,
	}

	if err := u.userService.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "User is", user)
}
