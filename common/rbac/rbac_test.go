package rbac

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yinloo-ola/tt-app/common/rbac/models"
	"github.com/yinloo-ola/tt-app/util"
	sqlitestore "github.com/yinloo-ola/tt-app/util/store/sqlite-store"
)

func TestParallelInsert(t *testing.T) {
	path := "rbac_parallel.db"
	t.Cleanup(func() {
		errRemove := os.Remove(path)
		if errRemove != nil {
			t.Fatalf("fail to clean up rbac.db. please clean up manually")
		}
		_ = os.Remove(path + "-shm")
		_ = os.Remove(path + "-wal")
	})

	assert := assert.New(t)

	permissionStore, err := sqlitestore.NewStore[models.Permission](path)
	util.PanicErr(err)
	roleStore, err := sqlitestore.NewStore[models.Role](path)
	util.PanicErr(err)
	userStore, err := sqlitestore.NewStore[models.User](path)
	util.PanicErr(err)
	rbac := NewRbac(
		permissionStore, roleStore, userStore,
	)
	defer func() {
		errClose := rbac.Close()
		util.PanicErr(errClose)
	}()
	permsIn := make([]models.Permission, 0, 100)
	permChan := make(chan models.Permission, 100)
	go func(ch chan models.Permission) {
		for p := range ch {
			permsIn = append(permsIn, p)
		}
	}(permChan)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		i := i
		wg.Add(1)
		go func(wg1 *sync.WaitGroup) {
			defer wg1.Done()
			perm := models.Permission{
				Name:        fmt.Sprintf("Name %d", i),
				Description: fmt.Sprintf("Desc %d", i),
			}
			id, err := rbac.PermissionStore.Insert(perm)
			util.PanicErr(err)
			perm.ID = id
			permChan <- perm
		}(&wg)
	}
	wg.Wait()
	close(permChan)

	perms, err := rbac.PermissionStore.FindWhere()
	util.PanicErr(err)
	assert.Len(perms, 100)
	assert.ElementsMatch(permsIn, perms)
}

func TestRbac_HasPermission(t *testing.T) {
	path := "rbac_has_permission.db"
	t.Cleanup(func() {
		errRemove := os.Remove(path)
		if errRemove != nil {
			t.Fatalf("fail to clean up rbac.db. please clean up manually")
		}
		_ = os.Remove(path + "-shm")
		_ = os.Remove(path + "-wal")
	})

	assert := assert.New(t)

	permissionStore, err := sqlitestore.NewStore[models.Permission](path)
	util.PanicErr(err)
	roleStore, err := sqlitestore.NewStore[models.Role](path)
	util.PanicErr(err)
	userStore, err := sqlitestore.NewStore[models.User](path)
	util.PanicErr(err)
	rbac := NewRbac(
		permissionStore, roleStore, userStore,
	)
	defer func() {
		errClose := rbac.Close()
		util.PanicErr(errClose)
	}()
	permsIn := make([]models.Permission, 0, 100)
	permChan := make(chan models.Permission, 100)
	go func(ch chan models.Permission) {
		for p := range ch {
			permsIn = append(permsIn, p)
		}
	}(permChan)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		i := i
		wg.Add(1)
		go func(wg1 *sync.WaitGroup) {
			defer wg1.Done()
			perm := models.Permission{
				Name:        fmt.Sprintf("Name %d", i),
				Description: fmt.Sprintf("Desc %d", i),
			}
			id, err := rbac.PermissionStore.Insert(perm)
			util.PanicErr(err)
			perm.ID = id
			permChan <- perm
		}(&wg)
	}
	wg.Wait()
	time.Sleep(10 * time.Millisecond)
	close(permChan)

	roles := make([]models.Role, 0, 10)
	partSize := len(permsIn) / 10
	for i := 0; i < 10; i++ {
		startIndex := i * partSize
		endIndex := (i + 1) * partSize

		role := models.Role{
			Name:        fmt.Sprintf("role %d", i),
			Description: fmt.Sprintf("desc %d", i),
		}
		for j := startIndex; j < endIndex; j++ {
			role.Permissions = append(role.Permissions, permsIn[j].ID)
		}
		id, err := rbac.RoleStore.Insert(role)
		util.PanicErr(err)
		role.ID = id
		roles = append(roles, role)
	}

	users := make([]models.User, 0, 5)
	for i := 0; i < 5; i++ {
		start := i * 2
		end := (i + 1) * 2
		user := models.User{
			UserID: fmt.Sprintf("userid %d", i),
		}
		for j := start; j < end; j++ {
			user.Roles = append(user.Roles, roles[j].ID)
		}
		id, err := rbac.UserStore.Insert(user)
		util.PanicErr(err)
		user.ID = id
		users = append(users, user)
	}

	for i := 0; i < 5; i++ {
		start := i * 20
		end := (i + 1) * 20

		for j := start; j < end; j++ {
			hasPerm, err := rbac.HasPermission(users[i].UserID, permsIn[j].ID)
			util.PanicErr(err)
			assert.True(hasPerm)
		}
	}

	for i := 0; i < 5; i++ {
		start := i * 20
		end := (i + 1) * 20

		for j := 0; j < 100; j++ {
			if j >= start && j < end {
				continue
			}
			hasPerm, err := rbac.HasPermission(users[i].UserID, permsIn[j].ID)
			util.PanicErr(err)
			assert.False(hasPerm)
		}
	}

}

func TestRbac_GetUserPermissions(t *testing.T) {
	path := "rbac_has_permission.db"
	t.Cleanup(func() {
		errRemove := os.Remove(path)
		if errRemove != nil {
			t.Fatalf("fail to clean up rbac.db. please clean up manually")
		}
		_ = os.Remove(path + "-shm")
		_ = os.Remove(path + "-wal")
	})

	assert := assert.New(t)

	permissionStore, err := sqlitestore.NewStore[models.Permission](path)
	util.PanicErr(err)
	roleStore, err := sqlitestore.NewStore[models.Role](path)
	util.PanicErr(err)
	userStore, err := sqlitestore.NewStore[models.User](path)
	util.PanicErr(err)
	rbac := NewRbac(
		permissionStore, roleStore, userStore,
	)
	defer func() {
		errClose := rbac.Close()
		util.PanicErr(errClose)
	}()
	permsIn := make([]models.Permission, 0, 100)
	permChan := make(chan models.Permission, 100)
	go func(ch chan models.Permission) {
		for p := range ch {
			permsIn = append(permsIn, p)
		}
	}(permChan)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		i := i
		wg.Add(1)
		go func(wg1 *sync.WaitGroup) {
			defer wg1.Done()
			perm := models.Permission{
				Name:        fmt.Sprintf("Name %d", i),
				Description: fmt.Sprintf("Desc %d", i),
			}
			id, err := rbac.PermissionStore.Insert(perm)
			util.PanicErr(err)
			perm.ID = id
			permChan <- perm
		}(&wg)
	}
	wg.Wait()
	time.Sleep(10 * time.Millisecond)
	close(permChan)

	roles := make([]models.Role, 0, 10)
	partSize := len(permsIn) / 10
	for i := 0; i < 10; i++ {
		startIndex := i * partSize
		endIndex := (i + 1) * partSize

		role := models.Role{
			Name:        fmt.Sprintf("role %d", i),
			Description: fmt.Sprintf("desc %d", i),
		}
		for j := startIndex; j < endIndex; j++ {
			role.Permissions = append(role.Permissions, permsIn[j].ID)
		}
		id, err := rbac.RoleStore.Insert(role)
		util.PanicErr(err)
		role.ID = id
		roles = append(roles, role)
	}

	users := make([]models.User, 0, 5)
	for i := 0; i < 5; i++ {
		start := i * 2
		end := (i + 1) * 2
		user := models.User{
			UserID: fmt.Sprintf("userid %d", i),
		}
		for j := start; j < end; j++ {
			user.Roles = append(user.Roles, roles[j].ID)
		}
		id, err := rbac.UserStore.Insert(user)
		util.PanicErr(err)
		user.ID = id
		users = append(users, user)
	}

	for i := 0; i < 5; i++ {
		start := i * 20
		end := (i + 1) * 20

		perms, err := rbac.GetUserPermissions(users[i].UserID)
		util.PanicErr(err)
		assert.ElementsMatch(perms, permsIn[start:end])
	}

}
