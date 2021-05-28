package utils

import (
	"crypto/rand"
	"fmt"
	"html/template"

	"github.com/russross/blackfriday"
)

func GenerateId() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
func ConvertMarkdownToHtml(markdown string) (html string) {
	return string(blackfriday.MarkdownBasic([]byte(markdown)))
}
func Unescape(x string) interface{} {
	return template.HTML(x)
}
