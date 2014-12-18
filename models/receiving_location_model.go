package models

type ReceivingLocation struct {
	Id              int    `db:"receiving_location_id" json:"receiving_location_id"`
	Type            string `db:"receiving_location_type" json:"receiving_location_type"`
	TemperatureZone string `db:"temperature_zone" json:"temperature_zone"`
}
