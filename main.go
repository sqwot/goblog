package main

import (
	"fmt"
	"goblog/models"
	"net/http"

	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
)

var posts map[string]*models.Post
var counter int

func indexHandler(rnd render.Render) {
	fmt.Println(counter)
	rnd.HTML(200, "index", posts)
}
func writeHandler(rnd render.Render) {
	rnd.HTML(200, "write", nil)
}
func savePostHandler(rnd render.Render, r *http.Request) {
	id := r.FormValue("id")
	title := r.FormValue("title")
	content := r.FormValue("content")

	var post *models.Post
	if id != "" {
		post = posts[id]
		post.Title = title
		post.Content = content
	} else {
		id = GenerateId()
		post := models.NewPost(id, title, content)
		posts[post.Id] = post
	}

	rnd.Redirect("/")
}

func editHandler(rnd render.Render, r *http.Request, params martini.Params) {
	id := params["id"]
	post, found := posts[id]
	if !found {
		rnd.Status(404)
		return
	}

	rnd.HTML(200, "write", post)
}

func deleteHandler(rnd render.Render, r *http.Request, params martini.Params) {
	id := params["id"]
	if id == "" {
		rnd.Status(404)
	}

	delete(posts, id)

	rnd.Redirect("/")
}

func main() {
	colorReset := "\033[0m"
	colorGreen := "\033[32m"
	fmt.Println(string(colorGreen), "***************BugaBlog on martini***************")
	fmt.Println(string(colorReset), "")

	posts = make(map[string]*models.Post, 0)
	counter = 0
	m := martini.Classic()
	m.Use(render.Renderer(render.Options{
		Directory:  "templates",
		Layout:     "layout",
		Extensions: []string{".tmpl", ".html"},
		//Funcs: ,
		Charset:    "UTF-8",
		IndentJSON: true,
	}))

	staticOptions := martini.StaticOptions{Prefix: "assets"}
	m.Use(martini.Static("assets", staticOptions))

	m.Get("/", indexHandler)
	m.Get("/write", writeHandler)
	m.Get("/edit/:id", editHandler)
	m.Get("/DeletePost/:id", deleteHandler)
	m.Post("/SavePost", savePostHandler)

	m.Run()
}
