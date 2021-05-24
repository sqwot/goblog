package main

import (
	"fmt"
	"html/template"
	"net/http"

	"goblog/models"
)

var posts map[string]*models.Post

func indexHandler(rw http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(rw, err.Error())
	}

	t.ExecuteTemplate(rw, "index", posts)
}
func writeHandler(rw http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/write.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(rw, err.Error())
	}

	t.ExecuteTemplate(rw, "write", nil)
}
func savePostHandler(w http.ResponseWriter, r *http.Request) {
	id := GenerateId()
	title := r.FormValue("title")
	content := r.FormValue("content")

	post := models.NewPost(id, title, content)
	posts[post.Id] = post

	http.Redirect(w, r, "/", 302)
}

func main() {
	colorReset := "\033[0m"
	colorGreen := "\033[32m"
	fmt.Println(string(colorGreen), "***************BugaBlog***************")
	fmt.Println(string(colorReset), "")

	posts = make(map[string]*models.Post, 0)
	fmt.Println("Application started at port 3000")

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/write", writeHandler)
	http.HandleFunc("/SavePost", savePostHandler)

	http.ListenAndServe(":3000", nil)
}
