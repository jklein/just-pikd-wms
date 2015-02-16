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

func TestSearchReceivingByZoneAndType(t *testing.T) {
	H(t).Test("GET", "/locations/receiving?temperature_zone=dry&type=DSD+Receiving+Bay", "").Check().HasStatus(200).BodyContains(`"rcl_id": "204-178900331-6"`)
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

/*
	router.HandleFunc("/locations/pickup", lc.Action(lc.CreatePickupLocation)).Methods("POST")
	router.HandleFunc("/locations/pickup/{id:[0-9]+}", lc.Action(lc.GetPickupLocation)).Methods("GET")
	router.HandleFunc("/locations/pickup/{id:[0-9]+}", lc.Action(lc.UpdatePickupLocation)).Methods("PATCH")
*/
func TestGetPickContainer(t *testing.T) {
	H(t).Test("GET", "/locations/containers/204-268900325-3", "").Check().HasStatus(200).BodyContains(`"pc_temperature_zone": "cold"`)
}

func TestGetPickContainerNotFound(t *testing.T) {
	H(t).Test("GET", "/locations/containers/204-168900281-5", "").Check().HasStatus(404)
}

func TestPatchPickContainer(t *testing.T) {
	body := `{"pc_id": "204-268900325-3", "pc_pcl_id": "204-268900322-2"}`
	H(t).Test("PATCH", "/locations/containers/204-268900325-3", body).Check().HasStatus(204)
	H(t).Test("GET", "/locations/containers/204-268900325-3", "").Check().HasStatus(200).BodyContains(`"pc_pcl_id": "204-268900322-2"`)
}

func TestPatchPickContainerNotFound(t *testing.T) {
	body := `{"pc_id": "204-168900281-5", "pc_pcl_id": "204-268900322-2"}`
	H(t).Test("PATCH", "/locations/containers/204-168900281-5", body).Check().HasStatus(404)
}

func TestPatchPickContainerIdentifierMismatch(t *testing.T) {
	body := `{"pc_id": "204-178900281-4", "pc_pcl_id": "204-268900322-2"}`
	H(t).Test("PATCH", "/locations/containers/204-268900325-3", body).Check().HasStatus(400)
}

func TestPostPickContainer(t *testing.T) {
	body := `{"pc_id": "204-168900324-9", "pc_pcl_id": null, "pc_temperature_zone": "cold",
    	"pc_type": "Bin", "pc_height": 18, "pc_width": 12, "pc_depth": 18}`
	H(t).Test("POST", "/locations/containers", body).Check().HasStatus(201)
	//make sure it got there
	H(t).Test("GET", "/locations/containers/204-168900324-9", "").Check().HasStatus(200).BodyContains(`"pc_temperature_zone": "cold"`)
}

func TestGetPickContainerLocation(t *testing.T) {
	H(t).Test("GET", "/locations/container_locations/204-268900322-2", "").Check().HasStatus(200).BodyContains(`"pcl_temperature_zone": "cold"`)
}

func TestGetPickContainerLocationNotFound(t *testing.T) {
	H(t).Test("GET", "/locations/container_locations/204-178900282-1", "").Check().HasStatus(404)
}

func TestPatchPickContainerLocation(t *testing.T) {
	body := `{"pcl_id": "204-268900322-2", "pcl_type": "Pick Cart Parking"}`
	H(t).Test("PATCH", "/locations/container_locations/204-268900322-2", body).Check().HasStatus(204)
	H(t).Test("GET", "/locations/container_locations/204-268900322-2", "").Check().HasStatus(200).BodyContains(`"pcl_type": "Pick Cart Parking"`)
}

func TestPatchPickContainerLocationNotFound(t *testing.T) {
	body := `{"pcl_id": "204-178900282-1", "pcl_type": "Pick Cart Parking"}`
	H(t).Test("PATCH", "/locations/container_locations/204-178900282-1", body).Check().HasStatus(404)
}

func TestPatchPickContainerLocationIdentifierMismatch(t *testing.T) {
	body := `{"pcl_id": "204-268900322-2", "pcl_type": "Pick Cart Parking"}`
	H(t).Test("PATCH", "/locations/container_locations/204-268900323-9", body).Check().HasStatus(400)
}

func TestPostPickContainerLocation(t *testing.T) {
	body := `{"pcl_id": "204-168900281-5", "pcl_type": "Finished Goods Buffer", "pcl_temperature_zone": "cold", "pcl_aisle": 1, "pcl_bay": 1, "pcl_shelf": 1, "pcl_shelf_slot": 1}`
	H(t).Test("POST", "/locations/container_locations", body).Check().HasStatus(201)
	//make sure it got there
	H(t).Test("GET", "/locations/container_locations/204-168900281-5", "").Check().HasStatus(200).BodyContains(`"pcl_temperature_zone": "cold"`)
}

func TestGetPickupLocation(t *testing.T) {
	H(t).Test("GET", "/locations/pickup/1", "").Check().HasStatus(200).BodyContains(`"pul_display_name": "Lane 1"`)
}

func TestGetPickupLocationNotFound(t *testing.T) {
	H(t).Test("GET", "/locations/pickup/1000", "").Check().HasStatus(404)
}

func TestPatchPickupLocation(t *testing.T) {
	body := `{"pul_id": 1, "pul_display_name": "Test Name Change"}`
	H(t).Test("PATCH", "/locations/pickup/1", body).Check().HasStatus(204)
	H(t).Test("GET", "/locations/pickup/1", "").Check().HasStatus(200).BodyContains(`"pul_display_name": "Test Name Change"`)
}

func TestPatchPickupLocationNotFound(t *testing.T) {
	body := `{"pul_id": 1000, "pul_display_name": "Test Name Change"}`
	H(t).Test("PATCH", "/locations/pickup/1000", body).Check().HasStatus(404)
}

func TestPatchPickupLocationIdentifierMismatch(t *testing.T) {
	body := `{"pul_id": 1000, "pul_display_name": "Test Name Change"}`
	H(t).Test("PATCH", "/locations/pickup/1", body).Check().HasStatus(400)
}

func TestPostPickupLocation(t *testing.T) {
	body := `{"pul_type": "Indoor Pickup Location", "pul_display_name": "Indoor 2", "pul_current_cars": 0}`
	H(t).Test("POST", "/locations/pickup", body).Check().HasStatus(201)
	//make sure it got there
	H(t).Test("GET", "/locations/pickup/12", "").Check().HasStatus(200).BodyContains(`"pul_display_name": "Indoor 2"`)
}
