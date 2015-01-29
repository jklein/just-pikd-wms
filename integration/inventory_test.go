// Copyright G2G Market Inc, 2015

package integration

import (
	"testing"
)

func TestGetStatic(t *testing.T) {
	H(t).Test("GET", "/inventory/static/1", "").Check().HasStatus(200).BodyContains(`"thumbnail_url": "https://s3.amazonaws.com/g2gcdn/68/00016000147720_200x200.jpg"`)
}

func TestGetStaticNotFound(t *testing.T) {
	H(t).Test("GET", "/inventory/static/10000", "").Check().HasStatus(404)
}

func TestPostStatic(t *testing.T) {
	body := `{"si_stl_id": "204-168900281-5", "si_pr_sku": "001-600014772-0"}`
	H(t).Test("POST", "/inventory/static", body).Check().HasStatus(201).BodyContains(`"si_pr_sku": "001-600014772-0"`)
}

func TestPatchStatic(t *testing.T) {
	body := `{"si_id": 1, "si_emptied_date": "2015-01-23T19:24:34.163739Z", "si_available_qty": 0}`
	H(t).Test("PATCH", "/inventory/static/1", body).Check().HasStatus(204)
	H(t).Test("GET", "/inventory/static/1", "").Check().HasStatus(200).BodyContains(`"si_emptied_date": "2015-01-23T19:24:34.163739Z"`)
}

func TestPatchStaticIdMismatch(t *testing.T) {
	body := `{"si_id": 2, "si_emptied_date": "2015-01-23T19:24:34.163739Z", "si_available_qty": 0}`
	H(t).Test("PATCH", "/inventory/static/1", body).Check().HasStatus(400)
}

func TestPatchStaticNotFound(t *testing.T) {
	body := `{"si_id": 1000, "si_emptied_date": "2015-01-23T19:24:34.163739Z", "si_available_qty": 0}`
	H(t).Test("PATCH", "/inventory/static/1000", body).Check().HasStatus(404)
}
