package routes

import (
	"GoBlog/db/documents"
	"GoBlog/models"
	"GoBlog/session"
	"fmt"

	"github.com/martini-contrib/render"
	"gopkg.in/mgo.v2"
)

func IndexHandler(rnd render.Render, s *session.Session, db *mgo.Database) {
	postCollection := db.C("posts")

	fmt.Println(s.Username)

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
