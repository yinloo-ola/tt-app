package api

import (
	"github.com/gin-gonic/gin"
)

func AddAPIs(routerGroup *gin.RouterGroup) {
	ctrl := &APIHomeController{}
	routerGroup.GET("/:link", ctrl.Link)
	routerGroup.GET("/", ctrl.Index)
}

type APIHomeController struct {
}

func (o *APIHomeController) Index(c *gin.Context) {
	c.HTML(200, "homeHTML", map[string]string{})
}

func (o *APIHomeController) Link(c *gin.Context) {
	link := c.Param("link")
	isHx := c.GetHeader("HX-Request")
	if isHx == "true" {
		c.HTML(200, link, map[string]string{})
		return
	}
	c.HTML(200, link+"HTML", map[string]string{})
}
