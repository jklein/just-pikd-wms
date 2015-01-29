// Copyright G2G Market Inc, 2015

package integration

import (
	"testing"
)

func TestGetStockingLocation(t *testing.T) {
	H(t).Test("GET", "/locations/stocking/204-168900281-5", "").Check().HasStatus(200).BodyContains(`"stl_temperature_zone": "dry"`)
}

func TestGetStockingLocationInvalid(t *testing.T) {
	H(t).Test("GET", "/locations/stocking/1", "").Check().HasStatus(500).BodyContains("invalid input syntax for EAN13")
}

func TestPostStockingLocation(t *testing.T) {
	body := `{"stl_id": "204-178900281-4", "stl_temperature_zone": "dry", "stl_type": "Pallet Storage"}`
	H(t).Test("POST", "/locations/stocking", body).Check().HasStatus(201)
}

func TestPatchStockingLocation(t *testing.T) {
	body := `{"stl_id": "204-168900281-5", "stl_needs_qc": true}`
	H(t).Test("PATCH", "/locations/stocking/204-168900281-5", body).Check().HasStatus(204)
	H(t).Test("GET", "/locations/stocking/204-168900281-5", "").Check().HasStatus(200).BodyContains(`"stl_needs_qc": true`)
}

func TestPatchStockingLocationIdentifierMismatch(t *testing.T) {
	body := `{"stl_id": "204-178900281-4", "stl_needs_qc": true}`
	H(t).Test("PATCH", "/locations/stocking/204-168900281-5", body).Check().HasStatus(400)
}

func TestPatchNonexistentStockingLocation(t *testing.T) {
	body := `{"stl_id": "204-178900282-1", "stl_needs_qc": true}`
	H(t).Test("PATCH", "/locations/stocking/204-178900282-1", body).Check().HasStatus(404)
}

func TestGetReceiving(t *testing.T) {
	H(t).Test("GET", "/locations/receiving/204-178900331-6", "").Check().HasStatus(200).BodyContains(`"rcl_temperature_zone": "dry"`)
}

func TestGetReceivingNotFound(t *testing.T) {
	H(t).Test("GET", "/locations/receiving/204-168900281-5", "").Check().HasStatus(404)
}

func TestSearchReceiving(t *testing.T) {
	H(t).Test("GET", "/locations/receiving?temperature_zone=dry", "").Check().HasStatus(200).BodyContains(`"rcl_id": "204-178900281-4"`)
}

func TestPatchAndSearchReceiving(t *testing.T) {
	body := `{"rcl_id": "204-178900284-5", "rcl_shi_shipment_code": 14}`
	H(t).Test("PATCH", "/locations/receiving/204-178900284-5", body).Check().HasStatus(204)
	H(t).Test("GET", "/locations/receiving/204-178900284-5", "").Check().HasStatus(200).BodyContains(`"rcl_shi_shipment_code": 14`)
	H(t).Test("GET", "/locations/receiving?temperature_zone=dry&has_product=true", "").Check().HasStatus(200).BodyContains(`"rcl_id": "204-178900284-5`)
}

func TestPatchReceivingNotFound(t *testing.T) {
	body := `{"rcl_id": "204-168900281-5", "rcl_shi_shipment_code": 14}`
	H(t).Test("PATCH", "/locations/receiving/204-168900281-5", body).Check().HasStatus(404)
}

func TestPatchReceivingIdentifierMismatch(t *testing.T) {
	body := `{"rcl_id": "204-178900281-4", "rcl_shi_shipment_code": 14}`
	H(t).Test("PATCH", "/locations/receiving/204-178900284-5", body).Check().HasStatus(400)
}

func TestPostReceiving(t *testing.T) {
	body := `{"rcl_id": "204-168900324-9", "rcl_type": "Pallet Receiving", "rcl_temperature_zone": "dry", "rcl_shi_shipment_code": null}`
	H(t).Test("POST", "/locations/receiving", body).Check().HasStatus(201)
	//make sure it got there
	H(t).Test("GET", "/locations/receiving/204-168900324-9", "").Check().HasStatus(200).BodyContains(`"rcl_temperature_zone": "dry"`)
}
