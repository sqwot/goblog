package models

type Post struct {
	Id              string
	Title           string
	ContentHtml     string
	ContentMarkdown string
}

func NewPost(id, title, contentHtml, contentMarkdown string) *Post {
	return &Post{ Id: id, Title: title, ContentHtml: contentHtml, ContentMarkdown:contentMarkdown }
}






