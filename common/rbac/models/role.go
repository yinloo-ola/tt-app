package models

import (
	"encoding/json"

	"github.com/yinloo-ola/tt-app/util"
	"github.com/yinloo-ola/tt-app/util/store"
)

type Role struct {
	ID          int64   `db:"id,pk" form:"id"`
	Name        string  `db:"name" form:"name"`
	Description string  `db:"description" form:"description"`
	Permissions []int64 `db:"permissions,json" form:"permissions"`
}

func (o *Role) FieldsVals() []any {
	perms, err := json.Marshal(o.Permissions)
	util.PanicErr(err)
	return []any{o.ID, o.Name, o.Description, perms}
}

func (o *Role) ScanRow(row store.RowScanner) error {
	var perms []byte
	err := row.Scan(&o.ID, &o.Name, &o.Description, &perms)
	if err != nil {
		return err
	}
	err = json.Unmarshal(perms, &o.Permissions)
	util.PanicErr(err)
	return nil
}
