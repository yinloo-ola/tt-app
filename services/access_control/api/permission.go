package api

import (
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

func (o *APIAccessController) PermissionModal(ctx *gin.Context) {
	slog.Debug("PermissionModal")
	permissionIDStr, _ := ctx.GetQuery("id")
	permissionID, err := strconv.ParseInt(permissionIDStr, 10, 64)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("invalid id"))
		return
	}
	actionType, _ := ctx.GetQuery("actionType")
	slog.Debug("PermissionModal", "id", permissionID, "actionType", actionType)
	permission, err := o.RbacStore.PermissionStore.GetOne(permissionID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("permission not found: %d", permissionID))
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to retrieve permission: %d. Error: %v", permissionID, err))
		return
	}

	ctx.HTML(200, "modal_once", gin.H{
		"IsHidden":  false,
		"ElementID": "update-permission-modal",
		"Body": o.templates.TemplateHTML("permission_modal", gin.H{
			"Action":      "update",
			"Name":        permission.Name,
			"Description": permission.Description,
			"ID":          permission.ID,
		}),
	})
}

func (o *APIAccessController) GetPermissions(ctx *gin.Context) {
	slog.Debug("GetPermissions")
	permissions, err := o.RbacStore.PermissionStore.FindWhere()
	if err != nil {
		slog.ErrorContext(ctx, "RbacStore.PermissionStore.FindWhere()", slog.String("error", err.Error()))
		_ = ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("fail to retrieve permissions"))
		return
	}

	permissionsContent := gin.H{
		"Permissions": permissions,
		"NewPermissionModal": gin.H{
			"IsHidden":  true,
			"ElementID": "new-permission-modal",
			"Body": o.templates.TemplateHTML("permission_modal", gin.H{
				"Action": "new",
			}),
		},
	}

	isHx := ctx.GetHeader("HX-Request")
	if isHx == "true" {
		if ctx.GetHeader("Hx-Target") == "ac-contents" {
			ctx.HTML(200, "permissions", permissionsContent)
			return
		}
		ctx.HTML(200, "access_control", gin.H{
			"Body": o.templates.TemplateHTML("permissions", permissionsContent),
		})
		return
	}

	ctx.HTML(200, "base", gin.H{
		"Title": "TT App - Access Control",
		"App":   "Table Tennis App",
		"Main": o.templates.TemplateHTML("access_control", gin.H{
			"Body": o.templates.TemplateHTML("permissions", permissionsContent),
		}),
	})
}

func (o *APIAccessController) AddPermission(ctx *gin.Context) {
	slog.Debug("AddPermission")
	var permission models.Permission
	err := ctx.Bind(&permission)
	if err != nil {
		slog.ErrorContext(ctx, "ctx.Bind()", slog.String("error", err.Error()))
		_ = ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("fail to bind body to permission"))
		return
	}
	slog.Debug("permission to add", "permission", permission)
	id, err := o.RbacStore.PermissionStore.Insert(permission)
	if err != nil {
		slog.ErrorContext(ctx, "RbacStore.PermissionStore.Insert()", slog.String("error", err.Error()))
		if errors.Is(err, store.ErrConflicted) {

			ctx.HTML(409, "error", gin.H{
				"ElementID": "permission-form-error",
				"Body":      template.HTML("Permission with the same name exists"),
			})
			ctx.Header("HX-Retarget", "#permission-form-error")
			return
		}
		_ = ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("fail to insert permission"))
		return
	}
	permission.ID = id
	ctx.HTML(200, "permission_row", permission)
}

func (o *APIAccessController) UpdatePermission(ctx *gin.Context) {
	var permission models.Permission
	err := ctx.Bind(&permission)
	if err != nil {
		slog.ErrorContext(ctx, "ctx.BindJSON()", slog.String("error", err.Error()))
		_ = ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("fail to bind body to permission"))
		return
	}
	slog.Debug("update permission", "permission", permission)
	err = o.RbacStore.PermissionStore.Update(permission.ID, permission)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			slog.ErrorContext(ctx, "RbacStore.PermissionStore.Update()", slog.String("error", err.Error()))
			_ = ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("permission not found"))
			return
		}
		if errors.Is(err, store.ErrConflicted) {
			ctx.HTML(409, "error", gin.H{
				"ElementID": "permission-form-error",
				"Body":      template.HTML("Permission with the same name exists"),
			})
			return
		}
		slog.ErrorContext(ctx, "RbacStore.PermissionStore.Update()", slog.String("error", err.Error()))
		_ = ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("fail to insert permission"))
		return
	}
	ctx.HTML(200, "permission_row", permission)
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
