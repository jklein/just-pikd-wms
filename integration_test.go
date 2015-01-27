// Copyright G2G Market Inc, 2015

package main

import (
	"fmt"
	"just-pikd-wms/config"
	. "just-pikd-wms/helpers"
	"net/http"
	"os"
	"testing"
)

// This file runs integration tests, which use the entire service end-to-end including
// the real webserver and database

// TestMain handles test setup and then runs tests
func TestMain(m *testing.M) {
	//setup -- TODO DRY this with main.go
	// initialize and load config
	config := &config.Config{}
	config.Load()

	// create DB connection pool
	db := SetupDB(config.DbUser, config.DbPass, config.DbName)

	// set up routes
	router := MakeRouter(db, config)

	// create negroni middleware handler
	n := MakeNegroni(router, config)
	//run the app in a background goroutine while testing continues
	go n.Run("localhost:3002")

	//start by invoking the handler to reset data, and make sure that succeeds
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://localhost:3002/reset", nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error resetting data: %s", err.Error())
		os.Exit(1)
	} else if resp.StatusCode != 200 {
		fmt.Printf("Error resetting data: status %d", resp.StatusCode)
		os.Exit(1)
	}

	//once data has been reset, we're ready to run tests
	status := m.Run()
	//tear down
	os.Exit(status)
}

//maybe it does make sense to put these in the controllers package?
//can also check coverage for other packages...
//could potentially make a package containing the server/routing logic, so that tests inside the controllers package can access that and get
//at the handler. however that still does not solve the data setup problem.

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

//func TestPatchSPOProduct(t *testing.T)

/*

#should succeed
curl -4 -i -X PATCH http://localhost:3000/spos/1 -d '{
    "spo_id": 1,
    "products": [
      {
        "spop_id": 109,
        "spop_confirmed_qty": 2,
        "spop_received_qty": 2
      }
    ]
}'

#should succeed and update both the SPO and two products
curl -4 -i -X PATCH http://localhost:3000/spos/1 -d '{
    "spo_id": 1,
    "spo_date_arrived": "2015-01-16T00:00:00Z",
    "products": [
      {
        "spop_id": 109,
        "spop_confirmed_qty": 3,
        "spop_received_qty": 3
      },
      {
        "spop_id": 111,
        "spop_confirmed_qty": 5,
        "spop_received_qty": 5
      }
    ]
}'

#should 404
curl -4 -i -X PATCH http://localhost:3000/spos/1 -d '{
    "spo_id": 1,
    "spo_date_arrived": "2015-01-14T00:00:00Z",
    "products": [
      {
        "spop_id": 123123123,
        "spop_confirmed_qty": 3,
        "spop_received_qty": 3
      },
      {
        "spop_id": 111,
        "spop_confirmed_qty": 2,
        "spop_received_qty": 2
      }
    ]
}'

#leaving the ID out of the embedded document entirely should also 404
curl -4 -i -X PATCH http://localhost:3000/spos/1 -d '{
    "spo_id": 1,
    "spo_date_arrived": "2015-01-14T00:00:00Z",
    "products": [
      {
        "spop_confirmed_qty": 3,
        "spop_received_qty": 3
      }
    ]
}'
*/
