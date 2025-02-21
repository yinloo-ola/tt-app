package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/yinloo-ola/tt-app/common/rbac/models"
	"github.com/yinloo-ola/tt-app/views/templ/access_control"
	"github.com/yinloo-ola/tt-app/views/templ/access_control/user"
	"github.com/yinloo-ola/tt-app/views/templ/base"
	"github.com/yinloo-ola/tt-app/views/templ/widget"
)

func (o *APIAccessController) GetUsers(ctx *gin.Context) {
	slog.Debug("GetUsers")

	users, err := o.RbacStore.UserStore.FindWhere()
	if err != nil {
		slog.ErrorContext(ctx, "RbacStore.UserStore.FindWhere()", "error", err)
		_ = ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("fail to retrieve users"))
		return
	}

	rows := make([]templ.Component, 0, len(users))
	for _, usr := range users {
		row := user.UserRow(usr.ID, usr.UserID)
		rows = append(rows, row)
	}

	roles, err := o.RbacStore.RoleStore.FindWhere()
	if err != nil {
		slog.ErrorContext(ctx, "RbacStore.RoleStore.FindWhere()", "error", err)
		_ = ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("fail to retrieve roles"))
		return
	}

	usersComponent := user.Users(rows, widget.Modal(
		true,
		"new-user-modal",
		user.UserModal("new", user.UserForm("new", 0, "", roles, []models.Role{}))),
	)

	isHx := ctx.GetHeader("HX-Request")
	if isHx == "true" {
		if ctx.GetHeader("Hx-Target") == "ac-contents" {
			ctx.HTML(200, "", usersComponent)
			return
		}
		ctx.HTML(200, "", access_control.AccessControl(usersComponent, "user"))
		return
	}

	ctx.HTML(200, "", base.Base("TT App - Access Control", "Table Tennis App", access_control.AccessControl(usersComponent, "role")))
}
