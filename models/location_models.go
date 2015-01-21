// Copyright G2G Market Inc, 2015

package models

import (
	"fmt"
	"gopkg.in/guregu/null.v2"
	"time"
)

type ReceivingLocation struct {
	Id              string   `json:"rcl_id"`
	Type            string   `json:"rcl_type"`
	TemperatureZone string   `json:"rcl_temperature_zone"`
	ShiShipmentCode null.Int `json:"rcl_shi_shipment_code"`
}

type StockingLocation struct {
	Id              string      `json:"stl_id"`
	TemperatureZone string      `json:"stl_temperature_zone"`
	Type            string      `json:"stl_type"`
	PickSegment     int         `json:"stl_pick_segment"`
	Aisle           int         `json:"stl_aisle"`
	Bay             int         `json:"stl_bay"`
	Shelf           int         `json:"stl_shelf"`
	ShelfSlot       int         `json:"stl_shelf_slot"`
	Height          null.Float  `json:"stl_height"`
	Width           null.Float  `json:"stl_width"`
	Depth           null.Float  `json:"stl_depth"`
	AssignedSku     null.String `json:"stl_assigned_sku"`
	NeedsQc         null.Bool   `json:"stl_needs_qc"`
	LastQcDate      *time.Time  `json:"stl_last_qc_date"`
	LocationCode    string      `json:"location_code"`
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
