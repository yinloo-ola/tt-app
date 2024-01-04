package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yinloo-ola/tt-app/common/rbac"
	"github.com/yinloo-ola/tt-app/common/rbac/models"
	"github.com/yinloo-ola/tt-app/util"
	sqlitestore "github.com/yinloo-ola/tt-app/util/store/sqlite-store"
)

func AddAPIs(routerGroup *gin.RouterGroup) {
	path := "rbac.db"
	permissionStore, err := sqlitestore.NewStore[models.Permission](path)
	util.PanicErr(err)
	roleStore, err := sqlitestore.NewStore[models.Role](path)
	util.PanicErr(err)
	userStore, err := sqlitestore.NewStore[models.User](path)
	util.PanicErr(err)
	rbacStore := rbac.NewRbac(
		permissionStore, roleStore, userStore,
	)
	ctrl := &APIAccessController{
		RbacStore: rbacStore,
	}
	routerGroup.GET("/permissions", ctrl.GetPermissions)
	routerGroup.POST("/permissions", ctrl.AddPermission)
	routerGroup.PUT("/permissions", ctrl.UpdatePermission)
	routerGroup.GET("/permission_modal", ctrl.PermissionModal)
	routerGroup.DELETE("/permissions/:id", ctrl.DeletePermission)

	routerGroup.GET("/roles", ctrl.GetRoles)
	routerGroup.GET("/role_modal", ctrl.RoleModal)

	routerGroup.POST("/roles", ctrl.AddRole)
	routerGroup.PUT("/roles", ctrl.UpdateRole)
	routerGroup.DELETE("/roles/:id", ctrl.DeleteRole)
}

type APIAccessController struct {
	RbacStore *rbac.Rbac
}
