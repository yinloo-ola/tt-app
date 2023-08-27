package api

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yinloo-ola/tt-app/common/rbac/models"
	"github.com/yinloo-ola/tt-app/util/store"
)

func (o *APIAccessController) GetRoles(ctx *gin.Context) {
	slog.Debug("GetRoles")
	role, err := o.RbacStore.RoleStore.FindWhere()
	if err != nil {
		slog.ErrorContext(ctx, "RbacStore.RoleStore.FindWhere()", slog.String("error", err.Error()))
		_ = ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("fail to retrieve roles"))
		return
	}
	ctx.JSON(http.StatusOK, role)
}

func (o *APIAccessController) AddRole(ctx *gin.Context) {
	slog.Debug("AddRole")
	var role models.Role
	err := ctx.BindJSON(&role)
	if err != nil {
		slog.ErrorContext(ctx, "ctx.BindJSON()", slog.String("error", err.Error()))
		_ = ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("fail to bind body to role"))
		return
	}
	id, err := o.RbacStore.RoleStore.Insert(role)
	if err != nil {
		slog.ErrorContext(ctx, "RbacStore.RoleStore.Insert()", slog.String("error", err.Error()))
		_ = ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("fail to insert role"))
		return
	}
	role.Id = id
	ctx.JSON(http.StatusOK, role)
}

func (o *APIAccessController) UpdateRole(ctx *gin.Context) {
	var role models.Role
	err := ctx.BindJSON(&role)
	if err != nil {
		slog.ErrorContext(ctx, "ctx.BindJSON()", slog.String("error", err.Error()))
		_ = ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("fail to bind body to role"))
		return
	}
	err = o.RbacStore.RoleStore.Update(role.Id, role)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			slog.ErrorContext(ctx, "RbacStore.RoleStore.Insert()", slog.String("error", err.Error()))
			_ = ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("role not found"))
			return
		}
		slog.ErrorContext(ctx, "RbacStore.RoleStore.Insert()", slog.String("error", err.Error()))
		_ = ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("fail to insert role"))
		return
	}
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
