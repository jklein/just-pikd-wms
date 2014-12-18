package models

type Supplier struct {
	Id   int    `db:"supplier_id" json:"supplier_id"`
	Name string `db:"supplier_name" json:"supplier_name"`
}
