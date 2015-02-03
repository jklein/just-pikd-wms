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
	"strconv"
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

// GetReceivingLocation retrieves a stocking location based on its id
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

// GetStockingLocation retrieves a stocking location based on its id
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
	stl_id := mux.Vars(r)["id"]
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return err, http.StatusBadRequest
	}

	var location models.StockingLocation
	err = json.Unmarshal(body, &location)
	if err != nil {
		return err, http.StatusBadRequest
	} else if stl_id != location.Id {
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

// GetPickContainer retrieves a pick container based on its id
func (c *locationController) GetPickContainer(rw http.ResponseWriter, r *http.Request) (error, int) {
	pc_id := mux.Vars(r)["id"]

	pc, err := c.dao.GetPickContainer(pc_id)

	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	c.JSON(rw, http.StatusOK, pc)
	return nil, 0
}

// CreatePickContainer creates a new pick container based on a passed in JSON object
func (c *locationController) CreatePickContainer(rw http.ResponseWriter, r *http.Request) (error, int) {
	var pc models.PickContainer
	err := jsonDecode(r.Body, &pc)

	if err != nil {
		// return a 400 if the request body doesn't decode to a PickContainer
		return err, http.StatusBadRequest
	}

	// pass the decoded model to the dao to update the DB
	pc, err = c.dao.CreatePickContainer(pc)

	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	// no response needed since there are no auto-generated IDs
	return nil, http.StatusCreated
}

// UpdatePickContainer updates a pick container based on a passed in json
func (c *locationController) UpdatePickContainer(rw http.ResponseWriter, r *http.Request) (error, int) {
	pc_id := mux.Vars(r)["id"]
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return err, http.StatusBadRequest
	}

	var pc models.PickContainer
	err = json.Unmarshal(body, &pc)
	if err != nil {
		return err, http.StatusBadRequest
	} else if pc_id != pc.Id {
		return errors.New("Identifier does not match request body for pc_id"), http.StatusBadRequest
	}

	// also decode to a dict so that update statements can be handled
	dict, err := jsonToDict(body)
	if err != nil {
		return err, http.StatusBadRequest
	}

	// pass the decoded model to the dao to update the DB
	err = c.dao.UpdatePickContainer(pc, dict)

	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	return nil, http.StatusNoContent
}

// GetPickContainerLocation retrieves a pick container location based on its id
func (c *locationController) GetPickContainerLocation(rw http.ResponseWriter, r *http.Request) (error, int) {
	pcl_id := mux.Vars(r)["id"]

	pcl, err := c.dao.GetPickContainerLocation(pcl_id)

	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	c.JSON(rw, http.StatusOK, pcl)
	return nil, 0
}

// CreatePickContainerLocation creates a new pick container location based on a passed in JSON object
func (c *locationController) CreatePickContainerLocation(rw http.ResponseWriter, r *http.Request) (error, int) {
	var pcl models.PickContainerLocation
	err := jsonDecode(r.Body, &pcl)

	if err != nil {
		// return a 400 if the request body doesn't decode to a PickContainerLocation
		return err, http.StatusBadRequest
	}

	// pass the decoded model to the dao to update the DB
	pcl, err = c.dao.CreatePickContainerLocation(pcl)

	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	// no response needed since there are no auto-generated IDs
	return nil, http.StatusCreated
}

// UpdatePickContainerLocation updates a pick container location based on a passed in json
func (c *locationController) UpdatePickContainerLocation(rw http.ResponseWriter, r *http.Request) (error, int) {
	// extract identifier from url - while we don't use this, it helps follow REST principles to have it in the URI
	// and could later be used for something like varnish cache invalidation
	pcl_id := mux.Vars(r)["id"]
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return err, http.StatusBadRequest
	}

	var pcl models.PickContainerLocation
	err = json.Unmarshal(body, &pcl)
	if err != nil {
		return err, http.StatusBadRequest
	} else if pcl_id != pcl.Id {
		return errors.New("Identifier does not match request body for pcl_id"), http.StatusBadRequest
	}

	// also decode to a dict so that update statements can be handled
	dict, err := jsonToDict(body)
	if err != nil {
		return err, http.StatusBadRequest
	}

	// pass the decoded model to the dao to update the DB
	err = c.dao.UpdatePickContainerLocation(pcl, dict)

	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	return nil, http.StatusNoContent
}

// GetPickupLocation retrieves a pickup location based on its id
func (c *locationController) GetPickupLocation(rw http.ResponseWriter, r *http.Request) (error, int) {
	pul_id := mux.Vars(r)["id"]

	pul, err := c.dao.GetPickupLocation(pul_id)

	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	c.JSON(rw, http.StatusOK, pul)
	return nil, 0
}

// CreatePickupLocation creates a new pickup location based on a passed in JSON object
func (c *locationController) CreatePickupLocation(rw http.ResponseWriter, r *http.Request) (error, int) {
	var pul models.PickupLocation
	err := jsonDecode(r.Body, &pul)

	if err != nil {
		// return a 400 if the request body doesn't decode to a PickupLocation
		return err, http.StatusBadRequest
	}

	// pass the decoded model to the dao to update the DB
	pul, err = c.dao.CreatePickupLocation(pul)

	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	// return the created model so that client can find out the auto generated id
	c.JSON(rw, http.StatusCreated, pul)
	return nil, 0
}

// UpdatePickupLocation updates a pickup location based on a passed in json
func (c *locationController) UpdatePickupLocation(rw http.ResponseWriter, r *http.Request) (error, int) {
	// extract identifier from url - while we don't use this, it helps follow REST principles to have it in the URI
	// and could later be used for something like varnish cache invalidation
	pul_id, _ := strconv.Atoi(mux.Vars(r)["id"])
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return err, http.StatusBadRequest
	}

	var pul models.PickupLocation
	err = json.Unmarshal(body, &pul)
	if err != nil {
		return err, http.StatusBadRequest
	} else if pul_id != pul.Id {
		return errors.New("Identifier does not match request body for pul_id"), http.StatusBadRequest
	}

	// also decode to a dict so that update statements can be handled
	dict, err := jsonToDict(body)
	if err != nil {
		return err, http.StatusBadRequest
	}

	// pass the decoded model to the dao to update the DB
	err = c.dao.UpdatePickupLocation(pul, dict)

	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	return nil, http.StatusNoContent
}
