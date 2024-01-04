package api

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strconv"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/yinloo-ola/tt-app/common/rbac/models"
	"github.com/yinloo-ola/tt-app/util/store"
	"github.com/yinloo-ola/tt-app/views/templ/access_control"
	"github.com/yinloo-ola/tt-app/views/templ/access_control/role"
	"github.com/yinloo-ola/tt-app/views/templ/base"
	"github.com/yinloo-ola/tt-app/views/templ/widget"
)

func (o *APIAccessController) RoleModal(ctx *gin.Context) {
	slog.Debug("RoleModal")
	roleIDStr, _ := ctx.GetQuery("id")
	roleID, err := strconv.ParseInt(roleIDStr, 10, 64)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("invalid id"))
		return
	}
	permissions, err := o.RbacStore.PermissionStore.FindWhere()
	if err != nil {
		slog.ErrorContext(ctx, "RbacStore.PermissionStore.FindWhere()", slog.String("error", err.Error()))
		_ = ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("fail to retrieve permissions"))
		return
	}
	actionType, _ := ctx.GetQuery("actionType")
	slog.Debug("RoleModal", "id", roleID, "actionType", actionType)
	rol, err := o.RbacStore.RoleStore.GetOne(roleID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("role not found: %d", roleID))
			return
		}
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("failed to retrieve role: %d. Error: %v", roleID, err))
		return
	}
	selectedPermissions := make([]models.Permission, 0, len(rol.Permissions))
	for _, pid := range rol.Permissions {
		idx := slices.IndexFunc(permissions, func(p models.Permission) bool {
			return p.ID == pid
		})
		selectedPermissions = append(selectedPermissions, permissions[idx])
	}

	ctx.HTML(200, "", widget.Modal(false, "update-role-modal",
		role.RoleModal("update", role.RoleForm("update", rol.ID, rol.Name, rol.Description, permissions, selectedPermissions))))
}

func (o *APIAccessController) GetRoles(ctx *gin.Context) {
	slog.Debug("GetRoles")
	roles, err := o.RbacStore.RoleStore.FindWhere()
	if err != nil {
		slog.ErrorContext(ctx, "RbacStore.RoleStore.FindWhere()", slog.String("error", err.Error()))
		_ = ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("fail to retrieve roles"))
		return
	}

	rows := make([]templ.Component, 0, len(roles))
	for _, rol := range roles {
		row := role.RoleRow(rol.ID, rol.Name, rol.Description)
		rows = append(rows, row)
	}

	permissions, err := o.RbacStore.PermissionStore.FindWhere()
	if err != nil {
		slog.ErrorContext(ctx, "RbacStore.PermissionStore.FindWhere()", slog.String("error", err.Error()))
		_ = ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("fail to retrieve permissions"))
		return
	}

	rolesComp := role.Roles(rows, widget.Modal(
		true,
		"new-role-modal",
		role.RoleModal("new", role.RoleForm("new", 0, "", "", permissions, []models.Permission{})),
	))

	isHx := ctx.GetHeader("HX-Request")
	if isHx == "true" {
		if ctx.GetHeader("Hx-Target") == "ac-contents" {
			ctx.HTML(200, "", rolesComp)
			return
		}
		ctx.HTML(200, "", access_control.AccessControl(rolesComp, "role"))
		return
	}

	ctx.HTML(200, "", base.Base("TT App - Access Control", "Table Tennis App", access_control.AccessControl(rolesComp, "role")))
}

func (o *APIAccessController) AddRole(ctx *gin.Context) {
	slog.Debug("AddRole")
	var rol models.Role
	err := ctx.Bind(&rol)
	if err != nil {
		slog.ErrorContext(ctx, "ctx.Bind()", slog.String("error", err.Error()))
		_ = ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("fail to bind body to role"))
		return
	}
	slog.Debug("role to add", "role", rol)
	id, err := o.RbacStore.RoleStore.Insert(rol)
	if err != nil {
		slog.ErrorContext(ctx, "RbacStore.RoleStore.Insert()", slog.String("error", err.Error()))
		if errors.Is(err, store.ErrConflicted) {
			ctx.Header("HX-Retarget", "#role-form-error")
			ctx.Header("HX-Reswap", "outerHTML transition:true")
			ctx.HTML(409, "", base.Error("role-form-error", "Role with the same name exists"))
			return
		}
		_ = ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("fail to insert role"))
		return
	}
	rol.ID = id
	ctx.HTML(200, "", role.RoleRow(rol.ID, rol.Name, rol.Description))
}

func (o *APIAccessController) UpdateRole(ctx *gin.Context) {
	var rol models.Role
	err := ctx.Bind(&rol)
	if err != nil {
		slog.ErrorContext(ctx, "ctx.BindJSON()", slog.String("error", err.Error()))
		_ = ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("fail to bind body to role"))
		return
	}
	slog.Debug("update role", "role", rol)
	err = o.RbacStore.RoleStore.Update(rol.ID, rol)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			slog.ErrorContext(ctx, "RbacStore.RoleStore.Update()", slog.String("error", err.Error()))
			_ = ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("role not found"))
			return
		}
		if errors.Is(err, store.ErrConflicted) {
			ctx.Header("HX-Retarget", "#role-form-error")
			ctx.Header("HX-Reswap", "outerHTML transition:true")
			ctx.HTML(409, "", base.Error("role-form-error", "Role with the same name exists"))
			return
		}
		slog.ErrorContext(ctx, "RbacStore.RoleStore.Update()", slog.String("error", err.Error()))
		_ = ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("fail to insert role"))
		return
	}
	ctx.HTML(200, "", role.RoleRow(rol.ID, rol.Name, rol.Description))
}

func (o *APIAccessController) DeleteRole(ctx *gin.Context) {
	slog.Debug("DeleteRole")
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.ErrorContext(ctx, "strconv.ParseInt", slog.String("id", idStr), slog.String("error", err.Error()))
		_ = ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("fail to delete role"))
		return
	}
	err = o.RbacStore.RoleStore.DeleteMulti([]int64{id})
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			slog.ErrorContext(ctx, "RbacStore.RoleStore.DeleteMulti()", slog.String("error", err.Error()))
			_ = ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("role not found"))
			return
		}
		slog.ErrorContext(ctx, "RbacStore.RoleStore.DeleteMulti()", slog.String("error", err.Error()))
		_ = ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("fail to delete role"))
		return
	}
}
