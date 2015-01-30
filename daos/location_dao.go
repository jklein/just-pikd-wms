// Copyright G2G Market Inc, 2015

package daos

import (
	"database/sql"
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
		`SELECT stl_id, stl_temperature_zone, stl_type, stl_pick_segment, stl_aisle, stl_bay, stl_shelf,
        stl_shelf_slot, stl_height, stl_width, stl_depth, stl_assigned_sku, stl_needs_qc, stl_last_qc_date
        FROM stocking_locations
        WHERE stl_id = $1;`, stl_id)
	return location, err
}

// CreateStockingLocation creates a new stocking location record based on the passed in model
// Although there is no auto-generated id column, the passed in model is returned unchanged to keep the interface the same
// as other Create* functions
func (dao *LocationDAO) CreateStockingLocation(stl models.StockingLocation) (models.StockingLocation, error) {
	_, err := dao.DB.NamedExec(
		`INSERT INTO stocking_locations (stl_id, stl_temperature_zone, stl_type, stl_pick_segment, stl_aisle, stl_bay,
        stl_shelf, stl_shelf_slot, stl_height, stl_width, stl_depth, stl_assigned_sku, stl_needs_qc, stl_last_qc_date)
        VALUES (:stl_id, :stl_temperature_zone, :stl_type, :stl_pick_segment, :stl_aisle, :stl_bay, :stl_shelf,
        :stl_shelf_slot, :stl_height, :stl_width, :stl_depth, :stl_assigned_sku, :stl_needs_qc, :stl_last_qc_date)
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

// GetReceivingLocations retrieves a slice of ReceivingLocation models
// and can be filtered by temperature zone, location type or which locations have product in them
func (dao *LocationDAO) GetReceivingLocations(temperature_zone string, has_product bool, location_type string) ([]models.ReceivingLocation, error) {
	var locations []models.ReceivingLocation

	args := struct {
		TemperatureZone string `json:"temperature_zone"`
		LocationType    string `json:"location_type"`
	}{temperature_zone, location_type}

	sql_string := `SELECT rcl_id, rcl_type,
        rcl_temperature_zone, rcl_shi_shipment_code
        FROM receiving_locations
        `

	// slice of where clause conditions based on whether params are set to their 0-value or not
	var conditions []string
	if has_product {
		conditions = append(conditions, "rcl_shi_shipment_code IS NOT NULL")
	}
	if len(temperature_zone) > 0 {
		conditions = append(conditions, "rcl_temperature_zone = :temperature_zone")
	}
	if len(location_type) > 0 {
		conditions = append(conditions, "rcl_type = :location_type")
	}

	sql_string += buildWhereFromConditions(conditions) + " ORDER BY rcl_id"

	rows, err := dao.DB.NamedQuery(sql_string, args)
	if err != nil {
		return locations, err
	}
	defer rows.Close()

	for rows.Next() {
		var l models.ReceivingLocation
		err = rows.StructScan(&l)
		if err != nil {
			return locations, err
		}
		locations = append(locations, l)
	}
	err = rows.Err()

	if err == nil && len(locations) == 0 {
		return locations, sql.ErrNoRows
	}

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
