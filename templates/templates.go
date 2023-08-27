package templates

import (
	"embed"
	"html/template"
)

//go:embed pages/*
var Files embed.FS

func Parse() *template.Template {
	return template.Must(template.ParseFS(Files, "pages/*"))
}
