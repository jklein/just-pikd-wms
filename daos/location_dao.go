package daos

import (
	"github.com/jmoiron/sqlx"
	"just-pikd-wms/models"
)

type Location_DAO struct {
	*sqlx.DB
}

func (dao *Location_DAO) Get_StockingLocation(stocking_location_id string) (models.StockingLocation, error) {
	var location models.StockingLocation

	err := dao.DB.Get(&location,
		`SELECT stocking_location_id, temperature_zone, stocking_location_type, pick_segment,
        aisle, bay, shelf, shelf_slot, height, width, depth, assigned_sku
        FROM stocking_locations
        WHERE stocking_location_id = $1;`, stocking_location_id)
	return location, err
}

func (dao *Location_DAO) Get_ReceivingLocations(temperature_zone string, has_product bool) ([]models.ReceivingLocation, error) {
	var locations []models.ReceivingLocation
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
