package models

type Post struct {
	Id      string
	Title   string
	Content string
}

func NewPost(id, titile, content string) *Post {
	return &Post{Id: id, Title: titile, Content: content}
}
