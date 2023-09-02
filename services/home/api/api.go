package api

import (
	"bytes"
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/yinloo-ola/tt-app/views"
)

func AddAPIs(routerGroup *gin.RouterGroup, templates views.TemplateExecutor) {
	ctrl := &APIHomeController{templates: templates}
	routerGroup.GET("/:link", ctrl.Link)
	routerGroup.GET("/", ctrl.Index)
}

type APIHomeController struct {
	templates views.TemplateExecutor
}

func (o *APIHomeController) Index(c *gin.Context) {
	buf := bytes.NewBufferString("")
	o.templates.ExecuteTemplate(buf, "home", nil)

	c.HTML(200, "base", map[string]any{
		"Title": "Table Tennis App",
		"App":   "Table Tennis App",
		"Main":  template.HTML(buf.String()),
	})
}

func (o *APIHomeController) Link(c *gin.Context) {
	link := c.Param("link")
	isHx := c.GetHeader("HX-Request")
	if isHx == "true" {
		c.HTML(200, link, map[string]string{})
		return
	}
	buf := bytes.NewBufferString("")
	o.templates.ExecuteTemplate(buf, link, nil)
	c.HTML(200, "base", map[string]any{
		"Title": "Table Tennis App",
		"App":   "Table Tennis App",
		"Main":  template.HTML(buf.String()),
	})
}
