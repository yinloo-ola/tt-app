package api

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yinloo-ola/tt-app/common/rbac/models"

	"github.com/yinloo-ola/tt-app/util/store"
)

func (o *APIAccessController) GetPermissions(ctx *gin.Context) {
	slog.Debug("GetPermissions")
	permissions, err := o.RbacStore.PermissionStore.FindWhere()
	if err != nil {
		slog.ErrorContext(ctx, "RbacStore.PermissionStore.FindWhere()", slog.String("error", err.Error()))
		_ = ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("fail to retrieve permissions"))
		return
	}

	isHx := ctx.GetHeader("HX-Request")
	if isHx == "true" {
		if ctx.GetHeader("Hx-Target") == "ac-contents" {
			ctx.HTML(200, "permissions", permissions)
			return
		}
		permissionsBuf := bytes.NewBufferString("")
		o.templates.ExecuteTemplate(permissionsBuf, "permissions", permissions)

		ctx.HTML(200, "access_control", map[string]any{
			"Body": template.HTML(permissionsBuf.String()),
		})
		return
	}
	permissionsBuf := bytes.NewBufferString("")
	o.templates.ExecuteTemplate(permissionsBuf, "permissions", permissions)
	buf := bytes.NewBufferString("")
	o.templates.ExecuteTemplate(buf, "access_control", map[string]any{
		"Body": template.HTML(permissionsBuf.String()),
	})
	ctx.HTML(200, "base", map[string]any{
		"Title": "TT App - Access Control",
		"App":   "Table Tennis App",
		"Main":  template.HTML(buf.String()),
	})
	// ctx.HTML(http.StatusOK, permission)
}

func (o *APIAccessController) AddPermission(ctx *gin.Context) {
	slog.Debug("AddPermission")
	var permission models.Permission
	err := ctx.BindJSON(&permission)
	if err != nil {
		slog.ErrorContext(ctx, "ctx.BindJSON()", slog.String("error", err.Error()))
		_ = ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("fail to bind body to permission"))
		return
	}
	id, err := o.RbacStore.PermissionStore.Insert(permission)
	if err != nil {
		slog.ErrorContext(ctx, "RbacStore.PermissionStore.Insert()", slog.String("error", err.Error()))
		_ = ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("fail to insert permission"))
		return
	}
	permission.Id = id
	ctx.JSON(http.StatusOK, permission)
}

func (o *APIAccessController) UpdatePermission(ctx *gin.Context) {
	var permission models.Permission
	err := ctx.BindJSON(&permission)
	if err != nil {
		slog.ErrorContext(ctx, "ctx.BindJSON()", slog.String("error", err.Error()))
		_ = ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("fail to bind body to permission"))
		return
	}
	err = o.RbacStore.PermissionStore.Update(permission.Id, permission)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			slog.ErrorContext(ctx, "RbacStore.PermissionStore.Insert()", slog.String("error", err.Error()))
			_ = ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("permission not found"))
			return
		}
		slog.ErrorContext(ctx, "RbacStore.PermissionStore.Insert()", slog.String("error", err.Error()))
		_ = ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("fail to insert permission"))
		return
	}
}

func (o *APIAccessController) DeletePermission(ctx *gin.Context) {
	slog.Debug("DeletePermission")
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.ErrorContext(ctx, "strconv.ParseInt", slog.String("id", idStr), slog.String("error", err.Error()))
		_ = ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("fail to delete permission"))
		return
	}
	err = o.RbacStore.PermissionStore.DeleteMulti([]int64{id})
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			slog.ErrorContext(ctx, "RbacStore.PermissionStore.DeleteMulti()", slog.String("error", err.Error()))
			_ = ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("permission not found"))
			return
		}
		slog.ErrorContext(ctx, "RbacStore.PermissionStore.DeleteMulti()", slog.String("error", err.Error()))
		_ = ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("fail to delete permission"))
		return
	}
}
