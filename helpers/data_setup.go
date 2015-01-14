// Copyright G2G Market Inc, 2015

// Package helpers contains helpers for testing
package helpers

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"just-pikd-wms/daos"
	"just-pikd-wms/models"
)

//DEV ONLY
//reset all test data
//panics if anything goes wrong
func Reset(DB *sqlx.DB) {
	DB.MustExec("TRUNCATE TABLE stocking_purchase_orders")
	DB.MustExec("ALTER SEQUENCE stocking_purchase_orders_stocking_purchase_order_id_seq RESTART")
	DB.MustExec("TRUNCATE TABLE stocking_purchase_order_products")
	DB.MustExec("ALTER SEQUENCE stocking_purchase_order_produ_stocking_purchase_order_produ_seq RESTART")
	DB.MustExec("TRUNCATE TABLE static_inventory")
	DB.MustExec("ALTER SEQUENCE static_inventory_static_inventory_id_seq RESTART")

	loadTestSPOs(DB)
	loadTestStatic(DB)
	return
}

func loadTestStatic(DB *sqlx.DB) {
	data, err := ioutil.ReadFile("./test_data/static.json")
	if err != nil {
		panic(err)
	}
	var sc []models.StaticInventory
	if err := json.Unmarshal(data, &sc); err != nil {
		panic(err)
	}

	dao := daos.InventoryDAO{DB: DB}
	for _, x := range sc {
		if _, err := dao.CreateStatic(x); err != nil {
			panic(err)
		}
	}
}

func loadTestSPOs(DB *sqlx.DB) {
	data, err := ioutil.ReadFile("./test_data/spos.json")
	if err != nil {
		panic(err)
	}
	var sc []models.StockingPurchaseOrder
	if err := json.Unmarshal(data, &sc); err != nil {
		panic(err)
	}

	dao := daos.StockingPurchaseOrderDAO{DB: DB}
	for _, x := range sc {
		if _, err := dao.CreateSPO(x); err != nil {
			panic(err)
		}
	}
}
