// Copyright G2G Market Inc, 2015

package daos

import (
	"github.com/jmoiron/sqlx"
	"just-pikd-wms/models"
)

// LocationDAO is used for data access related to
// locations such as stocking and receiving locations
type LocationDAO struct {
	*sqlx.DB
}

// GetStockingLocation retrieves a stocking location by its primary key, stl_id
func (dao *LocationDAO) GetStockingLocation(stl_id string) (models.StockingLocation, error) {
	var location models.StockingLocation

	err := dao.DB.Get(&location,
		`SELECT stl_id, stl_temperature_zone, stl_type, stl_pick_segment,
        stl_aisle, stl_bay, stl_shelf, stl_shelf_slot, stl_height, stl_width, stl_depth, stl_assigned_sku
        FROM stocking_locations
        WHERE stl_id = $1;`, stl_id)
	return location, err
}

// CreateStockingLocation creates a new stocking location record based on the passed in model
// Although there is no auto-generated id column, the passed in model is returned unchanged to keep the interface the same
// as other Create* functions
func (dao *LocationDAO) CreateStockingLocation(stl models.StockingLocation) (models.StockingLocation, error) {
	_, err := dao.DB.NamedExec(
		`INSERT INTO stocking_locations (stl_id, stl_temperature_zone, stl_type, stl_pick_segment,
        stl_aisle, stl_bay, stl_shelf, stl_shelf_slot, stl_height, stl_width, stl_depth, stl_assigned_sku)
        VALUES (:stl_id, :stl_temperature_zone, :stl_type, :stl_pick_segment,
        :stl_aisle, :stl_bay, :stl_shelf, :stl_shelf_slot, :stl_height, :stl_width, :stl_depth, :stl_assigned_sku)
        RETURNING stl_id`,
		stl)
	return stl, err
}

// UpdateStockingLocation updates a receiving location, updating only the passed-in fields
func (dao *LocationDAO) UpdateStockingLocation(location models.StockingLocation, dict map[string]interface{}) error {
	// update the base SPO object if any of its fields were passed in
	stmt := buildPatchUpdate("stocking_locations", "stl_id", dict)
	err := execCheckRows(dao.DB, stmt, location)
	return err
}

// GetReceivingLocation retrieves a receiving location by its primary key, rcl_id
func (dao *LocationDAO) GetReceivingLocation(rcl_id string) (models.ReceivingLocation, error) {
	var location models.ReceivingLocation

	err := dao.DB.Get(&location,
		`SELECT rcl_id, rcl_type,
        rcl_temperature_zone, rcl_shi_shipment_code
        FROM receiving_locations
        WHERE rcl_id = $1;`, rcl_id)
	return location, err
}

// GetReceivingLocations retrieves locations from a temperature zone
// It returns a slice of ReceivingLocation models for the temperature zone
// and can be filtered to retrieve only those locations that have product in them awaiting stocking
func (dao *LocationDAO) GetReceivingLocations(temperature_zone string, has_product bool) ([]models.ReceivingLocation, error) {
	var locations []models.ReceivingLocation

	//if has_product is true, add a where clause suffix looking for non-null shipment ids
	var where_suffix string
	if has_product {
		where_suffix = " AND rcl_shi_shipment_code IS NOT NULL"
	}

	sql := `SELECT rcl_id, rcl_type,
        rcl_temperature_zone, rcl_shi_shipment_code
        FROM receiving_locations
        WHERE rcl_temperature_zone = $1` + where_suffix + " ORDER BY rcl_id"

	err := dao.DB.Select(&locations, sql, temperature_zone)
	return locations, err
}

// UpdateReceivingLocation updates a receiving location, updating only the passed-in fields
func (dao *LocationDAO) UpdateReceivingLocation(location models.ReceivingLocation, dict map[string]interface{}) error {
	// update the base SPO object if any of its fields were passed in
	stmt := buildPatchUpdate("receiving_locations", "rcl_id", dict)
	err := execCheckRows(dao.DB, stmt, location)
	return err
}

// CreateReceivingLocation creates a new receiving location record based on the passed in model
func (dao *LocationDAO) CreateReceivingLocation(rcl models.ReceivingLocation) (models.ReceivingLocation, error) {
	_, err := dao.DB.NamedExec(
		`INSERT INTO receiving_locations (rcl_id, rcl_type, rcl_temperature_zone, rcl_shi_shipment_code)
        VALUES (:rcl_id, :rcl_type, :rcl_temperature_zone, :rcl_shi_shipment_code)
        RETURNING rcl_id`,
		rcl)
	return rcl, err
}
