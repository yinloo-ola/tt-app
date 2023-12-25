package api

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/yinloo-ola/tt-app/common/rbac/models"
	"github.com/yinloo-ola/tt-app/util/store"
	"github.com/yinloo-ola/tt-app/views/templ/access_control"
	"github.com/yinloo-ola/tt-app/views/templ/access_control/permission"
	"github.com/yinloo-ola/tt-app/views/templ/base"
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
	perm, err := o.RbacStore.PermissionStore.GetOne(permissionID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("permission not found: %d", permissionID))
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to retrieve permission: %d. Error: %v", permissionID, err))
		return
	}

	ctx.HTML(200, "", base.ModalOnce(false, "update-permission-modal",
		permission.PermissionModal("update", permission.PermissionForm("update", perm.ID, perm.Name, perm.Description))))
}

func (o *APIAccessController) GetPermissions(ctx *gin.Context) {
	slog.Debug("GetPermissions")
	permissions, err := o.RbacStore.PermissionStore.FindWhere()
	if err != nil {
		slog.ErrorContext(ctx, "RbacStore.PermissionStore.FindWhere()", slog.String("error", err.Error()))
		_ = ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("fail to retrieve permissions"))
		return
	}

	rows := make([]templ.Component, 0, len(permissions))
	for _, perm := range permissions {
		row := permission.PermissionRow(perm.ID, perm.Name, perm.Description)
		rows = append(rows, row)
	}

	permissionsComp := permission.Permissions(rows, base.Modal(
		true,
		"new-permission-modal",
		permission.PermissionModal("new", permission.PermissionForm("new", 0, "", ""))),
	)

	isHx := ctx.GetHeader("HX-Request")
	if isHx == "true" {
		if ctx.GetHeader("Hx-Target") == "ac-contents" {
			ctx.HTML(200, "", permissionsComp)
			return
		}
		ctx.HTML(200, "", access_control.AccessControl(permissionsComp))
		return
	}

	ctx.HTML(200, "", base.Base("TT App - Access Control", "Table Tennis App", access_control.AccessControl(permissionsComp)))
}

func (o *APIAccessController) AddPermission(ctx *gin.Context) {
	slog.Debug("AddPermission")
	var perm models.Permission
	err := ctx.Bind(&perm)
	if err != nil {
		slog.ErrorContext(ctx, "ctx.Bind()", slog.String("error", err.Error()))
		_ = ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("fail to bind body to permission"))
		return
	}
	slog.Debug("permission to add", "permission", perm)
	id, err := o.RbacStore.PermissionStore.Insert(perm)
	if err != nil {
		slog.ErrorContext(ctx, "RbacStore.PermissionStore.Insert()", slog.String("error", err.Error()))
		if errors.Is(err, store.ErrConflicted) {
			ctx.Header("HX-Retarget", "#permission-form-error")
			ctx.HTML(409, "", base.Error("permission-form-error", "Permission with the same name exists"))
			return
		}
		_ = ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("fail to insert permission"))
		return
	}
	perm.ID = id
	ctx.HTML(200, "", permission.PermissionRow(perm.ID, perm.Name, perm.Description))
}

func (o *APIAccessController) UpdatePermission(ctx *gin.Context) {
	var perm models.Permission
	err := ctx.Bind(&perm)
	if err != nil {
		slog.ErrorContext(ctx, "ctx.BindJSON()", slog.String("error", err.Error()))
		_ = ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("fail to bind body to permission"))
		return
	}
	slog.Debug("update permission", "permission", perm)
	err = o.RbacStore.PermissionStore.Update(perm.ID, perm)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			slog.ErrorContext(ctx, "RbacStore.PermissionStore.Update()", slog.String("error", err.Error()))
			_ = ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("permission not found"))
			return
		}
		if errors.Is(err, store.ErrConflicted) {
			ctx.Header("HX-Retarget", "#permission-form-error")
			ctx.HTML(409, "", base.Error("permission-form-error", "Permission with the same name exists"))
			return
		}
		slog.ErrorContext(ctx, "RbacStore.PermissionStore.Update()", slog.String("error", err.Error()))
		_ = ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("fail to insert permission"))
		return
	}
	ctx.HTML(200, "", permission.PermissionRow(perm.ID, perm.Name, perm.Description))
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
