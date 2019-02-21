package controllers

import (
	"net/http"

	"github.com/alexander-emelyanenko/go-web-server/views"
)

// Users struct describes our controller
type Users struct {
	NewView *views.View
}

// NewUsers method returns Users struct
func NewUsers() *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "views/users/new.gohtml"),
	}
}

// New method handles request
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}
