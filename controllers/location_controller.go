// Copyright G2G Market Inc, 2015

package controllers

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"gopkg.in/unrolled/render.v1"
	"just-pikd-wms/daos"
	"just-pikd-wms/models"
	"net/http"
)

// locationController contains methods related to WMS locations
type locationController struct {
	baseController
	*render.Render
	*sqlx.DB
	dao daos.LocationDAO
}

// NewLocationController acts as an initializer for locationController, setting
// its dao and returning the instance
func NewLocationController(rend *render.Render, db *sqlx.DB) *locationController {
	c := new(locationController)
	c.Render = rend
	c.DB = db
	c.dao = daos.LocationDAO{DB: db}
	return c
}

// GetStockingLocation retrieves a stocking location based on its id
// note, this is actually part of stocking not receiving
func (c *locationController) GetStockingLocation(rw http.ResponseWriter, r *http.Request) (error, int) {
	stocking_location_id := mux.Vars(r)["id"]

	location, err := c.dao.GetStockingLocation(stocking_location_id)

	if err == sql.ErrNoRows {
		return err, http.StatusNotFound
	} else if err != nil {
		c.LogError(err.Error())
		return err, http.StatusInternalServerError
	}

	// set computed field for Location Code on the location model
	location.SetLocationCode()

	c.JSON(rw, http.StatusOK, location)
	return nil, http.StatusOK
}

// GetReceivingLocations retrieves an array of locations for a temperature zone and can
// filter on whether they have product in them or not if desired.
// note, this is actually part of stocking not receiving
func (c *locationController) GetReceivingLocations(rw http.ResponseWriter, r *http.Request) (error, int) {
	hp := r.FormValue("has_product")
	temperature_zone := r.FormValue("temperature_zone")
	var has_product bool

	// accept string values 1 or true for has_product as true, other values considered false
	if hp == "1" || hp == "true" {
		has_product = true
	}

	// retrieve slice of locations from dao based on params
	locations, err := c.dao.GetReceivingLocations(temperature_zone, has_product)

	//no need to throw error if no rows are found as that is a normal case here
	if err != nil && err != sql.ErrNoRows {
		c.LogError(err.Error())
		return err, http.StatusInternalServerError
	}

	c.JSON(rw, http.StatusOK, locations)
	return nil, http.StatusOK
}

// UpdateReceivingLocation updates a receiving location based on a passed in model
// and is used to set or unset the supplier_shipment_id field to mark it as full or empty with product
// update receiving location to mark it as empty or filled with a specific supplier_shipment_id
func (c *locationController) UpdateReceivingLocation(rw http.ResponseWriter, r *http.Request) (error, int) {
	// extract identifier from url - while we don't use this, it helps follow REST principles to have it in the URI
	// and could later be used for something like varnish cache invalidation
	receiving_location_id := mux.Vars(r)["id"]
	// parse request body for a JSON receiving location model
	decoder := json.NewDecoder(r.Body)
	var location models.ReceivingLocation
	err := decoder.Decode(&location)
	if err != nil || receiving_location_id != location.Id {
		// return a 400 if the request body doesn't decode to a ReceivingLocation
		// or if the identifier doesn't match the request body's ID
		return err, http.StatusBadRequest
	}

	// pass the decoded model to the dao to update the DB
	err = c.dao.UpdateReceivingLocation(location)

	if err == sql.ErrNoRows {
		// return 404 if the row was not found
		return err, http.StatusNotFound
	} else if err != nil {
		c.LogError(err.Error())
		return err, http.StatusInternalServerError
	}

	// no response body needed for succesful update, just return 200
	return nil, http.StatusOK
}

// CreateReceivingLocation creates a new receiving location based on a passed in JSON object
func (c *locationController) CreateReceivingLocation(rw http.ResponseWriter, r *http.Request) (error, int) {
	// parse request body for a JSON receiving location model
	decoder := json.NewDecoder(r.Body)
	var rcl models.ReceivingLocation
	err := decoder.Decode(&rcl)

	if err != nil {
		// return a 400 if the request body doesn't decode to a ReceivingLocation
		return err, http.StatusBadRequest
	}

	// pass the decoded model to the dao to update the DB
	rcl, err = c.dao.CreateReceivingLocation(rcl)

	if err != nil {
		c.LogError(err.Error())
		return err, http.StatusInternalServerError
	}

	// no response needed since there are no auto-generated IDs
	return nil, http.StatusOK
}

// CreateStockingLocation creates a new stocking location based on a passed in JSON object
func (c *locationController) CreateStockingLocation(rw http.ResponseWriter, r *http.Request) (error, int) {
	// parse request body for a JSON stocking location model
	decoder := json.NewDecoder(r.Body)
	var stl models.StockingLocation
	err := decoder.Decode(&stl)

	if err != nil {
		// return a 400 if the request body doesn't decode to a StockingLocation
		return err, http.StatusBadRequest
	}

	// pass the decoded model to the dao to update the DB
	stl, err = c.dao.CreateStockingLocation(stl)

	if err != nil {
		c.LogError(err.Error())
		return err, http.StatusInternalServerError
	}

	// no response needed since there are no auto-generated IDs
	return nil, http.StatusOK
}
