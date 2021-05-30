package routes

import (
	"GoBlog/session"
	"fmt"
	"net/http"

	"github.com/martini-contrib/render"
)

func GetLoginHandler(rnd render.Render) {
	rnd.HTML(200, "login", nil)
}

func PostLoginHandler(rnd render.Render, r *http.Request, s *session.Session) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	fmt.Println(password)

	s.Username = username

	rnd.Redirect("/")
}
