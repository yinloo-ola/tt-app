package template

import (
	"bytes"
	"fmt"
	"html/template"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/yinloo-ola/tt-app/views"
)

type TemplateExecutor interface {
	ExecuteTemplate(wr io.Writer, name string, data interface{}) error
	TemplateHTML(name string, data interface{}) template.HTML
}

type DebugTemplateExecutor struct {
	Engine *gin.Engine
}

func (e *DebugTemplateExecutor) ExecuteTemplate(wr io.Writer, name string, data interface{}) error {
	t := views.ParseFiles()
	return t.ExecuteTemplate(wr, name, data)
}

func (e *DebugTemplateExecutor) TemplateHTML(name string, data interface{}) template.HTML {
	buf := bytes.NewBufferString("")
	err := e.ExecuteTemplate(buf, name, data)
	if err != nil {
		panic(fmt.Sprintf("fail to execute template: %s, data: %#v, error: %v", name, data, err))
	}
	return template.HTML(buf.String())
}

type ReleaseTemplateExecutor struct {
	Template *template.Template
}

func (e *ReleaseTemplateExecutor) ExecuteTemplate(wr io.Writer, name string, data interface{}) error {
	return e.Template.ExecuteTemplate(wr, name, data)
}

func (e *ReleaseTemplateExecutor) TemplateHTML(name string, data interface{}) template.HTML {
	buf := bytes.NewBufferString("")
	err := e.ExecuteTemplate(buf, name, data)
	if err != nil {
		panic(fmt.Sprintf("fail to execute template: %s, data: %#v, error: %v", name, data, err))
	}
	return template.HTML(buf.String())
}
