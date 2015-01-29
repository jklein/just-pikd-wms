// Copyright G2G Market Inc, 2015

package integration

import (
	"testing"
)

func TestGetShipment(t *testing.T) {
	H(t).Test("GET", "/suppliers/shipments/1", "").Check().HasStatus(200).BodyContains(`"shi_shipment_code": "425565063"`)
}

func TestGetShipmentNotFound(t *testing.T) {
	H(t).Test("GET", "/suppliers/shipments/10000", "").Check().HasStatus(404)
}

func TestSearchAllShipments(t *testing.T) {
	H(t).Test("GET", "/suppliers/shipments", "").Check().HasStatus(200).BodyContains(`"shi_shipment_code": "425565063"`)
}

func TestSearchShipmentsByCode(t *testing.T) {
	H(t).Test("GET", "/suppliers/shipments?shipment_code=777152188", "").Check().HasStatus(200).BodyContains(`"shi_id": 3`)
}

func TestSearchShipmentsBySPO(t *testing.T) {
	H(t).Test("GET", "/suppliers/shipments?spo_id=3", "").Check().HasStatus(200).BodyContains(`"shi_id": 3`)
}

func TestSearchShipmentsBySPOAndCode(t *testing.T) {
	H(t).Test("GET", "/suppliers/shipments?shipment_code=777152188&spo_id=3", "").Check().HasStatus(200).BodyContains(`"shi_id": 3`)
}

func TestSearchShipmentsBySPOAndCodeNotFound(t *testing.T) {
	H(t).Test("GET", "/suppliers/shipments?shipment_code=777152188&spo_id=1", "").Check().HasStatus(404)
}

func TestCreateShipment(t *testing.T) {
	body := `{"shi_shipment_code": "5923997121", "shi_spo_id": 1, "shi_su_id": 1}`
	H(t).Test("POST", "/suppliers/shipments", body).Check().HasStatus(201).BodyContains(`"shi_shipment_code": "5923997121"`)
}

func TestPatchShipment(t *testing.T) {
	body := `{"shi_id": 1, "shi_actual_delivery": "2015-01-17T00:00:00Z"}`
	H(t).Test("PATCH", "/suppliers/shipments/1", body).Check().HasStatus(204)
	H(t).Test("GET", "/suppliers/shipments/1", "").Check().HasStatus(200).BodyContains(`"shi_actual_delivery": "2015-01-17T00:00:00Z"`)
}

func TestPatchShipmentNotFound(t *testing.T) {
	body := `{"shi_id": 1000, "shi_actual_delivery": "2015-01-17T00:00:00Z"}`
	H(t).Test("PATCH", "/suppliers/shipments/1000", body).Check().HasStatus(404)
}

func TestPatchShipmentIdentifierMismatch(t *testing.T) {
	body := `{"shi_id": 1000, "shi_actual_delivery": "2015-01-17T00:00:00Z"}`
	H(t).Test("PATCH", "/suppliers/shipments/1", body).Check().HasStatus(400)
}
