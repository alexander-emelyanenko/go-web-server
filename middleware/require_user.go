package middleware

import (
	"fmt"
	"net/http"

	"github.com/alexander-emelyanenko/go-web-server/models"
)

type RequireUser struct {
	models.UserService
}

func (mv *RequireUser) ApplyFn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		user, err := mv.UserService.ByRemember(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
		}
		fmt.Println("User found: ", user)
		next(w, r)
	})
}

func (mv *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mv.ApplyFn(next.ServeHTTP)
}
