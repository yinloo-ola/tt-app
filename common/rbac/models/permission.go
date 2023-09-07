package models

import "github.com/yinloo-ola/tt-app/util/store"

type Permission struct {
	ID          int64  `db:"id,pk" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
}

func (o *Permission) FieldsVals() []any {
	return []any{o.ID, o.Name, o.Description}
}

func (o *Permission) ScanRow(row store.RowScanner) error {
	return row.Scan(&o.ID, &o.Name, &o.Description)
}
