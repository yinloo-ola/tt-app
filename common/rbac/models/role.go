package models

import (
	"encoding/json"

	"github.com/yinloo-ola/tt-app/util"
	"github.com/yinloo-ola/tt-app/util/store"
)

type Role struct {
	Id          int64   `db:"id,pk"`
	Name        string  `db:"name"`
	Description string  `db:"description"`
	Permissions []int64 `db:"permissions,json"`
}

func (o *Role) FieldsVals() []any {
	perms, err := json.Marshal(o.Permissions)
	util.PanicErr(err)
	return []any{o.Id, o.Name, o.Description, perms}
}

func (o *Role) ScanRow(row store.RowScanner) error {
	var perms []byte
	err := row.Scan(&o.Id, &o.Name, &o.Description, &perms)
	if err != nil {
		return err
	}
	err = json.Unmarshal(perms, &o.Permissions)
	util.PanicErr(err)
	return nil
}
