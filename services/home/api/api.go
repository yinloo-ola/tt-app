package api

import (
	"bytes"
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/yinloo-ola/tt-app/util/template_util"
)

func AddAPIs(routerGroup *gin.RouterGroup, templates template_util.TemplateExecutor) {
	ctrl := &APIHomeController{templates: templates}
	routerGroup.GET("/", ctrl.Index)
}

type APIHomeController struct {
	templates template_util.TemplateExecutor
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
