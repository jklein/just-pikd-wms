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

func (dao *LocationDAO) GetStockingLocation(stl_id string) (models.StockingLocation, error) {
	var location models.StockingLocation

	err := dao.DB.Get(&location,
		`SELECT stl_id, stl_temperature_zone, stl_type, stl_pick_segment,
        stl_aisle, stl_bay, stl_shelf, stl_shelf_slot, stl_height, stl_width, stl_depth, stl_assigned_sku
        FROM stocking_locations
        WHERE stl_id = $1;`, stl_id)
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
        WHERE rcl_temperature_zone = $1` + where_suffix

	err := dao.DB.Select(&locations, sql, temperature_zone)
	return locations, err
}

// PutReceivingLocation updates a receiving location and returns any errors it received
// The passed in model should be a location that already exists
// Only the mutable field - supplier_shipment_id, is updated to tag the location as in use or not in use
// Other fields are considered immutable and are not updated.
func (dao *LocationDAO) UpdateReceivingLocation(location models.ReceivingLocation) error {
	result, err := dao.DB.NamedExec(`UPDATE receiving_locations
        set rcl_shi_shipment_code = :rcl_shi_shipment_code,
        last_updated = now()
        WHERE rcl_id = :rcl_id`, location)

	if err == nil {
		// if the update doesn't match any rows, return this so the client knows it was unsuccessful
		if rows, _ := result.RowsAffected(); rows == 0 {
			return sql.ErrNoRows
		}
	}
	return err
}
