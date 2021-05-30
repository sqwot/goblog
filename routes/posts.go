package routes

import (
	"GoBlog/db/documents"
	"GoBlog/models"
	"GoBlog/session"
	"GoBlog/utils"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"gopkg.in/mgo.v2"
)

func WriteHandler(rnd render.Render, s *session.Session) {
	if !s.IsAuthorized {
		rnd.Redirect("/")
	}

	model := models.EditPostModel{}
	model.IsAuthorized = s.IsAuthorized
	model.Post = models.Post{}
	rnd.HTML(200, "write", model)
}
func SavePostHandler(rnd render.Render, r *http.Request, db *mgo.Database, s *session.Session) {
	if !s.IsAuthorized {
		rnd.Redirect("/")
	}

	postCollection := db.C("posts")
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

func EditHandler(rnd render.Render, r *http.Request, params martini.Params, db *mgo.Database, s *session.Session) {
	if !s.IsAuthorized {
		rnd.Redirect("/")
	}

	postCollection := db.C("posts")
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

	model := models.EditPostModel{}
	model.IsAuthorized = s.IsAuthorized
	model.Post = post

	rnd.HTML(200, "write", model)
}

func ViewHandler(rnd render.Render, r *http.Request, params martini.Params, db *mgo.Database, s *session.Session) {
	postCollection := db.C("posts")
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

	model := models.ViewPostModel{}
	model.IsAuthorized = s.IsAuthorized
	model.Post = post

	rnd.HTML(200, "view", model)
}

func DeleteHandler(rnd render.Render, r *http.Request, params martini.Params, db *mgo.Database, s *session.Session) {
	if !s.IsAuthorized {
		rnd.Redirect("/")
	}

	postCollection := db.C("posts")
	id := params["id"]
	if id == "" {
		rnd.Status(404)
		return
	}

	postCollection.RemoveId(id)

	rnd.Redirect("/")
}

func GetHtmlHandler(rnd render.Render, r *http.Request) {
	md := r.FormValue("md")
	html := utils.ConvertMarkdownToHtml(md)

	rnd.JSON(200, map[string]interface{}{"html": html})
}
