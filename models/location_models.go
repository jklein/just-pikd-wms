// Copyright G2G Market Inc, 2015

package models

import (
	"fmt"
	"gopkg.in/guregu/null.v2"
)

type ReceivingLocation struct {
	Id                 string   `db:"receiving_location_id" json:"receiving_location_id"`
	Type               string   `db:"receiving_location_type" json:"receiving_location_type"`
	TemperatureZone    string   `db:"temperature_zone" json:"temperature_zone"`
	SupplierShipmentId null.Int `db:"supplier_shipment_id" json:"supplier_shipment_id"`
}

type StockingLocation struct {
	Id                   string      `db:"stocking_location_id" json:"stocking_location_id"`
	TemperatureZone      string      `db:"temperature_zone" json:"temperature_zone"`
	StockingLocationType string      `db:"stocking_location_type" json:"stocking_location_type"`
	PickSegment          int         `db:"pick_segment" json:"pick_segment"`
	Aisle                int         `db:"aisle" json:"aisle"`
	Bay                  int         `db:"bay" json:"bay"`
	Shelf                int         `db:"shelf" json:"shelf"`
	ShelfSlot            int         `db:"shelf_slot" json:"shelf_slot"`
	Height               null.Float  `db:"height" json:"height"`
	Width                null.Float  `db:"width" json:"width"`
	Depth                null.Float  `db:"depth" json:"depth"`
	AssignedSku          null.String `db:"assigned_sku" json:"assigned_sku"`
	LocationCode         string      `json:"location_code"`
}

// SetLocationCode computes the location code value, which is a more human readable
// display of the location that can help to find it in the store.
func (loc *StockingLocation) SetLocationCode() {
	//short display letter for temperature zones
	var zone string
	switch loc.TemperatureZone {
	case "dry":
		zone = "D"
	case "frozen":
		zone = "F"
	case "cold":
		zone = "C"
	case "fresh", "perishable":
		zone = "P"
	default:
		zone = "U" //unknown
	}

	//result will look like D-A01-B01-S04-T03
	//%02 pads with leading zeros up to 2 total digits
	loc.LocationCode = fmt.Sprintf("%s-A%02d-B%02d-S%02d-T%02d",
		zone, loc.Aisle, loc.Bay, loc.Shelf, loc.ShelfSlot)

	return
}
