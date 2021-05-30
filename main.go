package main

import (
	"GoBlog/routes"
	"GoBlog/session"
	"GoBlog/utils"
	"fmt"
	"html/template"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"gopkg.in/mgo.v2"
)

func main() {
	colorReset := "\033[0m"
	colorGreen := "\033[32m"
	fmt.Println(string(colorGreen), "***************BugaBlog on martini***************")
	fmt.Println(string(colorReset), "")

	dbSession, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	db := dbSession.DB("BugaBlog")

	m := martini.Classic()

	m.Map(db)
	m.Use(session.Middleware)

	unescapeFuncMap := template.FuncMap{"unescape": utils.Unescape}

	m.Use(render.Renderer(render.Options{
		Directory:  "templates",
		Layout:     "layout",
		Extensions: []string{".tmpl", ".html"},
		Funcs:      []template.FuncMap{unescapeFuncMap},
		Charset:    "UTF-8",
		IndentJSON: true,
	}))

	staticOptions := martini.StaticOptions{Prefix: "assets"}
	m.Use(martini.Static("assets", staticOptions))

	m.Get("/", routes.IndexHandler)
	m.Get("/write", routes.WriteHandler)
	m.Get("/edit/:id", routes.EditHandler)
	m.Get("/view/:id", routes.ViewHandler)
	m.Get("/DeletePost/:id", routes.DeleteHandler)
	m.Get("/login", routes.GetLoginHandler)
	m.Get("/logout", routes.LogoutHandler)

	m.Post("/SavePost", routes.SavePostHandler)
	m.Post("/getHtml", routes.GetHtmlHandler)
	m.Post("/login", routes.PostLoginHandler)

	m.Run()
}
