// Copyright G2G Market Inc, 2015

package integration

import (
	//"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSPO(t *testing.T) {
	H(t).Test("GET", "/spos/1", "").Check().HasStatus(200).BodyContains(`"spo_id": 1`).BodyContains(`"spo_date_arrived": null`)
}

func TestGetSPONotFound(t *testing.T) {
	H(t).Test("GET", "/spos/1000000", "").Check().HasStatus(404)
}

func TestPatchSPO(t *testing.T) {
	H(t).Test("PATCH", "/spos/1", `{"spo_id": 1, "spo_date_arrived": "2015-01-14T00:00:00Z"}`).Check().HasStatus(200)
	H(t).Test("GET", "/spos/1", "").Check().HasStatus(200).BodyContains(`"spo_id": 1`).BodyContains(`"spo_date_arrived": "2015-01-14T00:00:00Z"`)
}

func TestPatchSPOBadRequest(t *testing.T) {
	H(t).Test("PATCH", "/spos/2", `{"spo_id": 1, "spo_date_arrived": "2015-01-14T00:00:00Z"}`).Check().HasStatus(400)
}

func TestPatchSPONotFound(t *testing.T) {
	H(t).Test("PATCH", "/spos/1000000", `{"spo_id": 1000000, "spo_date_arrived": "2015-01-14T00:00:00Z"}`).Check().HasStatus(404)
}

func TestPatchSPOProduct(t *testing.T) {
	body := `{"spo_id": 1,
        "products": [
            {"spop_id": 109, "spop_confirmed_qty": 2, "spop_received_qty": 2}
        ]}`
	H(t).Test("PATCH", "/spos/1", body).Check().HasStatus(200)
	//b := H(t).Test("GET", "/spos/1", "").Check().HasStatus(200).Body
	//TODO once there's an endpoint to get a single product, check that the response body has the modified fields, and still has other fields too
}

func TestPatchTwoSPOProducts(t *testing.T) {
	body := `{"spo_id": 1, "spo_date_arrived": "2015-01-16T00:00:00Z",
        "products": [
            {"spop_id": 109,"spop_confirmed_qty": 3,"spop_received_qty": 3},
            {"spop_id": 111,"spop_confirmed_qty": 5,"spop_received_qty": 5}
        ]}`
	H(t).Test("PATCH", "/spos/1", body).Check().HasStatus(200)
	//TODO verify
}

func TestPatchNonexistentProduct(t *testing.T) {
	body := `{"spo_id": 1, "spo_date_arrived": "2015-01-14T00:00:00Z",
        "products": [
            {"spop_id": 123123123, "spop_confirmed_qty": 3, "spop_received_qty": 3},
            {"spop_id": 111, "spop_confirmed_qty": 2, "spop_received_qty": 2}
        ]}`
	H(t).Test("PATCH", "/spos/1", body).Check().HasStatus(404)
}

func TestPatchProductMissingSpopId(t *testing.T) {
	body := `{"spo_id": 1, "spo_date_arrived": "2015-01-14T00:00:00Z",
    "products": [
        {"spop_confirmed_qty": 3, "spop_received_qty": 3},
    ]}`
	H(t).Test("PATCH", "/spos/1", body).Check().HasStatus(400)
}
