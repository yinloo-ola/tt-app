package views

import (
	"embed"
	"html/template"
)

//go:embed pages/**
var Files embed.FS

//go:embed assets/*
var Assets embed.FS

func ParseFS() *template.Template {
	return template.Must(template.ParseFS(Files, "pages/**/*"))
}
func ParseGlob() *template.Template {
	return template.Must(template.ParseGlob("views/pages/**/*"))
}
func Glob() string {
	return "views/pages/**/*"
}
