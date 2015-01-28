// Copyright G2G Market Inc, 2015

package daos

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"just-pikd-wms/daos"
	"just-pikd-wms/models"
	"reflect"
)

//DEV ONLY
//reset all test data
//panics if anything goes wrong
func ResetTestData(DB *sqlx.DB) {
	DB.MustExec("TRUNCATE TABLE stocking_purchase_orders")
	DB.MustExec("ALTER SEQUENCE stocking_purchase_orders_spo_id_seq RESTART")
	DB.MustExec("TRUNCATE TABLE stocking_purchase_order_products")
	DB.MustExec("ALTER SEQUENCE stocking_purchase_order_products_spop_id_seq RESTART")
	DB.MustExec("TRUNCATE TABLE static_inventory")
	DB.MustExec("ALTER SEQUENCE static_inventory_si_id_seq RESTART")
	DB.MustExec("TRUNCATE TABLE stocking_locations")
	DB.MustExec("TRUNCATE TABLE receiving_locations")
	DB.MustExec("TRUNCATE TABLE supplier_shipments")
	DB.MustExec("ALTER SEQUENCE supplier_shipments_shi_id_seq RESTART")

	loadTestSPOs(DB)
	loadTestStatic(DB)
	loadTestShipments(DB)
	loadTestStockingLocations(DB)
	loadTestReceivingLocations(DB)
	return
}

func loadTestStatic(DB *sqlx.DB) {
	data := loadFromFile("./test_data/static.json")
	dao := daos.InventoryDAO{DB: DB}
	var sc []models.StaticInventory
	if err := json.Unmarshal(data, &sc); err != nil {
		panic(err)
	}
	loadFromSlice(interfaceSlice(sc), dao.CreateStatic)
}

func loadTestSPOs(DB *sqlx.DB) {
	data := loadFromFile("./test_data/spos.json")
	dao := daos.StockingPurchaseOrderDAO{DB: DB}
	var sc []models.StockingPurchaseOrder
	if err := json.Unmarshal(data, &sc); err != nil {
		panic(err)
	}
	loadFromSlice(interfaceSlice(sc), dao.CreateSPO)
}

func loadTestShipments(DB *sqlx.DB) {
	data := loadFromFile("./test_data/supplier_shipments.json")
	dao := daos.SupplierDAO{DB: DB}
	var sc []models.SupplierShipment
	if err := json.Unmarshal(data, &sc); err != nil {
		panic(err)
	}
	loadFromSlice(interfaceSlice(sc), dao.CreateShipment)
}

func loadTestStockingLocations(DB *sqlx.DB) {
	data := loadFromFile("./test_data/stocking_locations.json")
	dao := daos.LocationDAO{DB: DB}
	var sc []models.StockingLocation
	if err := json.Unmarshal(data, &sc); err != nil {
		panic(err)
	}
	loadFromSlice(interfaceSlice(sc), dao.CreateStockingLocation)
}

func loadTestReceivingLocations(DB *sqlx.DB) {
	data := loadFromFile("./test_data/receiving_locations.json")
	dao := daos.LocationDAO{DB: DB}
	var sc []models.ReceivingLocation
	if err := json.Unmarshal(data, &sc); err != nil {
		panic(err)
	}
	loadFromSlice(interfaceSlice(sc), dao.CreateReceivingLocation)
}

// loadFromFile loads data from a file and panics if there are any errors
func loadFromFile(file_path string) []byte {
	data, err := ioutil.ReadFile(file_path)
	if err != nil {
		panic(err)
	}
	return data
}

// loadFromSlice is a generic method caller to save some verbosity on the loadX methods iterating over slice elements to call a function
func loadFromSlice(model_slice []interface{}, f interface{}) {
	fn := reflect.ValueOf(f)

	for _, x := range model_slice {
		if res := fn.Call([]reflect.Value{reflect.ValueOf(x)}); res[1].Interface() != nil {
			panic(fmt.Sprintf("%v loading from slice %v", res[1].Interface(), res[0].Interface()))
		}
	}
}

// interfaceSlice converts a passed in slice to a slice of interface values, suitable to be passed into a more generic method like loadFromSlice
func interfaceSlice(slice interface{}) []interface{} {
	v := reflect.ValueOf(slice)
	if v.Kind() != reflect.Slice {
		panic(fmt.Errorf("interfaceSlice: not a slice but %T", slice))
	}
	result := make([]interface{}, v.Len())
	for i := range result {
		result[i] = v.Index(i).Interface()
	}
	return result
}
