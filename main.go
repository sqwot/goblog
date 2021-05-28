package main

import (
	"GoBlog/db/documents"
	"GoBlog/session"
	"GoBlog/utils"
	"fmt"
	"goblog/models"
	"html/template"
	"net/http"
	"time"

	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
	"gopkg.in/mgo.v2"
)

const (
	COOKIE_NAME = "sessionId"
)

var postCollection *mgo.Collection
var inMemorySession *session.Session

func indexHandler(rnd render.Render, r *http.Request) {
	cookie, err := r.Cookie(COOKIE_NAME)
	if err == nil {
		fmt.Println(inMemorySession.Get(cookie.Value))
	}

	postDocuments := []documents.PostDocument{}
	postCollection.Find(nil).All(&postDocuments)

	posts := []models.Post{}

	for _, doc := range postDocuments {
		post := models.Post{
			Id:              doc.Id,
			Title:           doc.Title,
			ContentHtml:     doc.ContentHtml,
			ContentMarkdown: doc.ContentMarkdown,
		}
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
	contentHtml := utils.ConvertMarkdownToHtml(contentMarkdown)

	postDocument := documents.PostDocument{
		Id:              id,
		Title:           title,
		ContentHtml:     contentHtml,
		ContentMarkdown: contentMarkdown,
	}
	if id != "" {
		err := postCollection.UpdateId(id, postDocument)
		if err != nil {
			rnd.Status(401)
			return
		}
	} else {
		id = utils.GenerateId()
		postDocument.Id = id
		postCollection.Insert(postDocument)
	}

	rnd.Redirect("/")
}

func editHandler(rnd render.Render, r *http.Request, params martini.Params) {
	id := params["id"]

	postDocument := documents.PostDocument{}
	err := postCollection.FindId(id).One(&postDocument)
	if err != nil {
		rnd.Status(404)
		return
	}
	post := models.Post{
		Id:              postDocument.Id,
		Title:           postDocument.Title,
		ContentHtml:     postDocument.ContentHtml,
		ContentMarkdown: postDocument.ContentMarkdown,
	}

	rnd.HTML(200, "write", post)
}

func deleteHandler(rnd render.Render, r *http.Request, params martini.Params) {
	id := params["id"]
	if id == "" {
		rnd.Status(404)
		return
	}

	postCollection.RemoveId(id)

	rnd.Redirect("/")
}

func getHtmlHandler(rnd render.Render, r *http.Request) {
	md := r.FormValue("md")
	html := utils.ConvertMarkdownToHtml(md)

	rnd.JSON(200, map[string]interface{}{"html": html})
}

func getLoginHandler(rnd render.Render) {
	rnd.HTML(200, "login", nil)
}

func postLoginHandler(rnd render.Render, rw http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	fmt.Println(password)

	sessionId := inMemorySession.Init(username)
	cookie := &http.Cookie{
		Name:    COOKIE_NAME,
		Value:   sessionId,
		Expires: time.Now().Add(5 * time.Minute),
	}

	http.SetCookie(rw, cookie)

	rnd.Redirect("/")
}

func main() {
	colorReset := "\033[0m"
	colorGreen := "\033[32m"
	fmt.Println(string(colorGreen), "***************BugaBlog on martini***************")
	fmt.Println(string(colorReset), "")

	dbSession, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	postCollection = dbSession.DB("BugaBlog").C("posts")

	m := martini.Classic()

	inMemorySession = session.NewSession()

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

	m.Get("/", indexHandler)
	m.Get("/write", writeHandler)
	m.Get("/edit/:id", editHandler)
	m.Get("/DeletePost/:id", deleteHandler)
	m.Get("/login", getLoginHandler)

	m.Post("/SavePost", savePostHandler)
	m.Post("/getHtml", getHtmlHandler)
	m.Post("/login", postLoginHandler)

	m.Run()
}
