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
	DB.MustExec("TRUNCATE TABLE inbound_inventory")
	DB.MustExec("ALTER SEQUENCE inbound_inventory_inbound_inventory_id_seq RESTART")
	DB.MustExec("TRUNCATE TABLE stocking_purchase_orders")
	DB.MustExec("ALTER SEQUENCE stocking_purchase_orders_stocking_purchase_order_id_seq RESTART")
	DB.MustExec("TRUNCATE TABLE stocking_purchase_order_products")
	DB.MustExec("ALTER SEQUENCE stocking_purchase_order_produ_stocking_purchase_order_produ_seq RESTART")
	DB.MustExec("TRUNCATE TABLE static_inventory")
	DB.MustExec("ALTER SEQUENCE static_inventory_static_inventory_id_seq RESTART")

	loadTestInbound(DB)
	loadTestSPOs(DB)
	loadTestStatic(DB)
	return
}

func loadTestInbound(DB *sqlx.DB) {
	data, err := ioutil.ReadFile("./test_data/inbound.json")
	if err != nil {
		panic(err)
	}
	var ic []models.InboundInventory
	if err := json.Unmarshal(data, &ic); err != nil {
		panic(err)
	}

	dao := daos.Inventory_DAO{DB: DB}
	for _, x := range ic {
		if _, err := dao.Create_Inbound(x); err != nil {
			panic(err)
		}
	}
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

	dao := daos.Inventory_DAO{DB: DB}
	for _, x := range sc {
		if _, err := dao.Create_Static(x); err != nil {
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

	dao := daos.StockingPurchaseOrder_DAO{DB: DB}
	for _, x := range sc {
		if _, err := dao.Create_SPO(x); err != nil {
			panic(err)
		}
	}
}
