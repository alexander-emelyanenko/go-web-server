package controllers

import (
	"fmt"
	"net/http"

	"github.com/alexander-emelyanenko/go-web-server/views"
)

// NewUsers method returns Users struct
func NewUsers() *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "users/new"),
	}
}

// Users struct describes our users controller
type Users struct {
	NewView *views.View
}

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

	fmt.Fprintln(w, "Name is", form.Name)
	fmt.Fprintln(w, "Email is", form.Email)
	fmt.Fprintln(w, "Password is", form.Password)
}
