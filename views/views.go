package views

import (
	"embed"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

//go:embed pages/**
var Files embed.FS

//go:embed assets/*
var Assets embed.FS

func ParseFS() *template.Template {
	t := template.Must(template.ParseFS(Files, "pages/**/*.html"))
	return template.Must(t.ParseFS(Files, "pages/**/**/*.html"))
}

func ParseFiles() *template.Template {
	return template.Must(template.ParseFiles(GetFiles()...))
}

func GetFiles() []string {
	files := []string{}
	err := filepath.Walk("views/pages", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // you can also return nil here if you want to skip files and directories that can not be accessed
		}
		if info.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil // ignore directories and other files
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic("ParseFiles failed: " + err.Error())
	}
	return files
}
