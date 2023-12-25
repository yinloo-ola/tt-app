package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yinloo-ola/tt-app/views/templ/base"
	"github.com/yinloo-ola/tt-app/views/templ/home"
)

func AddAPIs(routerGroup *gin.RouterGroup) {
	ctrl := &APIHomeController{}
	routerGroup.GET("/", ctrl.Index)
}

type APIHomeController struct {
}

func (o *APIHomeController) Index(c *gin.Context) {
	c.HTML(200, "", base.Base("Table Tennis App", "Table Tennis App", home.Home()))
}
