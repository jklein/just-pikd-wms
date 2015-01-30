// Copyright G2G Market Inc, 2015

package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"gopkg.in/unrolled/render.v1"
	"io/ioutil"
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

	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	// set computed field for Location Code on the location model
	location.SetLocationCode()

	c.JSON(rw, http.StatusOK, location)
	return nil, 0
}

// GetReceivingLocation retrieves a stocking location based on its id
// note, this is actually part of stocking not receiving
func (c *locationController) GetReceivingLocation(rw http.ResponseWriter, r *http.Request) (error, int) {
	rcl_id := mux.Vars(r)["id"]

	location, err := c.dao.GetReceivingLocation(rcl_id)

	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	c.JSON(rw, http.StatusOK, location)
	return nil, 0
}

// GetReceivingLocations retrieves an array of locations for a temperature zone and can
// filter on whether they have product in them or not if desired.
func (c *locationController) GetReceivingLocations(rw http.ResponseWriter, r *http.Request) (error, int) {
	hp := r.FormValue("has_product")
	temperature_zone := r.FormValue("temperature_zone")
	location_type := r.FormValue("type")

	var has_product bool

	// accept string values 1 or true for has_product as true, other values considered false
	if hp == "1" || hp == "true" {
		has_product = true
	}

	// retrieve slice of locations from dao based on params
	locations, err := c.dao.GetReceivingLocations(temperature_zone, has_product, location_type)

	//no need to throw error if no rows are found as that is a normal case here
	if err != nil && err != sql.ErrNoRows {
		c.LogError(err.Error())
		return err, http.StatusInternalServerError
	}

	c.JSON(rw, http.StatusOK, locations)
	return nil, 0
}

// UpdateReceivingLocation updates a receiving location based on a passed in json
// and is used to set or unset the rcl_shi_id field to mark it as full or empty with product
// update receiving location to mark it as empty or filled with a specific supplier shipment
func (c *locationController) UpdateReceivingLocation(rw http.ResponseWriter, r *http.Request) (error, int) {
	// extract identifier from url - while we don't use this, it helps follow REST principles to have it in the URI
	// and could later be used for something like varnish cache invalidation
	receiving_location_id := mux.Vars(r)["id"]
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return err, http.StatusBadRequest
	}

	var location models.ReceivingLocation
	err = json.Unmarshal(body, &location)
	if err != nil {
		return err, http.StatusBadRequest
	} else if receiving_location_id != location.Id {
		return errors.New("Identifier does not match request body for rcl_id"), http.StatusBadRequest
	}

	// also decode to a dict so that update statements can be handled
	dict, err := jsonToDict(body)
	if err != nil {
		return err, http.StatusBadRequest
	}

	// pass the decoded model to the dao to update the DB
	err = c.dao.UpdateReceivingLocation(location, dict)

	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	return nil, http.StatusNoContent
}

// CreateReceivingLocation creates a new receiving location based on a passed in JSON object
func (c *locationController) CreateReceivingLocation(rw http.ResponseWriter, r *http.Request) (error, int) {
	var rcl models.ReceivingLocation
	err := jsonDecode(r.Body, &rcl)

	if err != nil {
		// return a 400 if the request body doesn't decode to a ReceivingLocation
		return err, http.StatusBadRequest
	}

	// pass the decoded model to the dao to update the DB
	rcl, err = c.dao.CreateReceivingLocation(rcl)

	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	// no response needed since there are no auto-generated IDs
	return nil, http.StatusCreated
}

// CreateStockingLocation creates a new stocking location based on a passed in JSON object
func (c *locationController) CreateStockingLocation(rw http.ResponseWriter, r *http.Request) (error, int) {
	var stl models.StockingLocation
	err := jsonDecode(r.Body, &stl)

	if err != nil {
		// return a 400 if the request body doesn't decode to a StockingLocation
		return err, http.StatusBadRequest
	}

	// pass the decoded model to the dao to update the DB
	stl, err = c.dao.CreateStockingLocation(stl)

	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	// no response needed since there are no auto-generated IDs
	return nil, http.StatusCreated
}

// UpdateStockingLocation updates a stocking location based on a passed in json
func (c *locationController) UpdateStockingLocation(rw http.ResponseWriter, r *http.Request) (error, int) {
	// extract identifier from url - while we don't use this, it helps follow REST principles to have it in the URI
	// and could later be used for something like varnish cache invalidation
	receiving_location_id := mux.Vars(r)["id"]
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return err, http.StatusBadRequest
	}

	var location models.StockingLocation
	err = json.Unmarshal(body, &location)
	if err != nil {
		return err, http.StatusBadRequest
	} else if receiving_location_id != location.Id {
		return errors.New("Identifier does not match request body for stl_id"), http.StatusBadRequest
	}

	// also decode to a dict so that update statements can be handled
	dict, err := jsonToDict(body)
	if err != nil {
		return err, http.StatusBadRequest
	}

	// pass the decoded model to the dao to update the DB
	err = c.dao.UpdateStockingLocation(location, dict)

	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	return nil, http.StatusNoContent
}
