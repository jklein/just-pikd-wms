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

func (dao *LocationDAO) GetStockingLocation(stocking_location_id string) (models.StockingLocation, error) {
	var location models.StockingLocation

	err := dao.DB.Get(&location,
		`SELECT stocking_location_id, temperature_zone, stocking_location_type, pick_segment,
        aisle, bay, shelf, shelf_slot, height, width, depth, assigned_sku
        FROM stocking_locations
        WHERE stocking_location_id = $1;`, stocking_location_id)
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
		where_suffix = " AND supplier_shipment_id IS NOT NULL"
	}

	sql := `SELECT receiving_location_id, receiving_location_type,
        temperature_zone, supplier_shipment_id
        FROM receiving_locations
        WHERE temperature_zone = $1` + where_suffix

	err := dao.DB.Select(&locations, sql, temperature_zone)
	return locations, err
}

// PutReceivingLocation updates a receiving location and returns any errors it received
// The passed in model should be a location that already exists
// Only the mutable field - supplier_shipment_id, is updated to tag the location as in use or not in use
// Other fields are considered immutable and are not updated.
func (dao *LocationDAO) PutReceivingLocation(location models.ReceivingLocation) error {
	result, err := dao.DB.NamedExec(`UPDATE receiving_locations
        set supplier_shipment_id = :supplier_shipment_id,
        last_updated = now()
        WHERE receiving_location_id = :receiving_location_id`, location)

	if err == nil {
		// if the update doesn't match any rows, return this so the client knows it was unsuccessful
		if rows, _ := result.RowsAffected(); rows == 0 {
			return sql.ErrNoRows
		}
	}
	return err
}
