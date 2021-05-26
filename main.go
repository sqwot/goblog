package main

import (
	"GoBlog/db/documents"
	"fmt"
	"goblog/models"
	"html/template"
	"net/http"

	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
	"gopkg.in/mgo.v2"
)

var postCollection *mgo.Collection

func indexHandler(rnd render.Render) {
	postDocuments := []documents.postDocument{}
	postCollection.Find(nil).All(&postDocuments)

	posts := []models.Post{}

	for _, doc := range postDocuments {
		post := models.Post{doc.Id, doc.Title, doc.ContentHtml, doc.ContentMarkdown}
		posts = append(posts, post)
	}

	rnd.HTML(200, "index", posts)
}
func writeHandler(rnd render.Render) {
	post := models.Post{}
	rnd.HTML(200, "write", post)
}
func savePostHandler(rnd render.Render, r *http.Request) {
	id := r.FormValue("id")
	title := r.FormValue("title")
	contentMarkdown := r.FormValue("content")
	contentHtml := ConvertMarkdownToHtml(contentMarkdown)

	var post *models.Post
	if id != "" {
		post = posts[id]
		post.Title = title
		post.ContentHtml = contentHtml
		post.ContentMarkdown = contentMarkdown
	} else {
		id = GenerateId()
		post := models.NewPost(id, title, contentHtml, contentMarkdown)
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

func getHtmlHandler(rnd render.Render, r *http.Request) {
	md := r.FormValue("md")
	html := ConvertMarkdownToHtml(md)

	rnd.JSON(200, map[string]interface{}{"html": html})
}

func unescape(x string) interface{} {
	return template.HTML(x)
}

func main() {
	colorReset := "\033[0m"
	colorGreen := "\033[32m"
	fmt.Println(string(colorGreen), "***************BugaBlog on martini***************")
	fmt.Println(string(colorReset), "")

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	postCollection = session.DB("BugaBlog").C("posts")

	m := martini.Classic()

	unescapeFuncMap := template.FuncMap{"unescape": unescape}

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

	m.Get("/", indexHandler)
	m.Get("/write", writeHandler)
	m.Get("/edit/:id", editHandler)
	m.Get("/DeletePost/:id", deleteHandler)
	m.Post("/SavePost", savePostHandler)
	m.Post("/getHtml", getHtmlHandler)

	m.Run()
}
