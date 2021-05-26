package documents

type postDocument struct {
	Id              string `bson:"_id,omitempty"`
	Title           string
	ContentHtml     string
	ContentMarkdown string
}
