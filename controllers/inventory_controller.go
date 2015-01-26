// Copyright G2G Market Inc, 2015

package controllers

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"gopkg.in/unrolled/render.v1"
	"just-pikd-wms/daos"
	"just-pikd-wms/helpers"
	"just-pikd-wms/models"
	"net/http"
	"strconv"
)

// inventoryController contains handlers for /inventory routes
type inventoryController struct {
	baseController
	*render.Render
	*sqlx.DB
	dao daos.InventoryDAO
}

// NewInventoryController acts as an initializer for inventoryController, setting
// its dao and returning the instance
func NewInventoryController(rend *render.Render, db *sqlx.DB) *inventoryController {
	c := new(inventoryController)
	c.Render = rend
	c.DB = db
	c.dao = daos.InventoryDAO{DB: db}
	return c
}

// Reset is a helper function for resetting test data in dev only
// It calls the helper to reload data from json files
func (c *inventoryController) Reset(rw http.ResponseWriter, r *http.Request) (error, int) {
	helpers.Reset(c.DB)
	rw.Write([]byte("Data reset complete"))
	return nil, http.StatusOK
}

// GetStatic retrieves a static inventory record by its id
func (c *inventoryController) GetStatic(rw http.ResponseWriter, r *http.Request) (error, int) {
	// parse args
	static_id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// retrieve static inventory model from dao by id
	static, err := c.dao.GetStatic(static_id)

	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	// set computed field on static inventory model for thumbnail image URL
	static.SetThumbnailURL()

	// render and return response
	c.JSON(rw, http.StatusOK, static)
	return nil, http.StatusOK
}

// CreateStatic creates a new static inventory record based on a passed in JSON object
func (c *inventoryController) CreateStatic(rw http.ResponseWriter, r *http.Request) (error, int) {
	var static models.StaticInventory
	err := jsonDecode(r.Body, &static)

	if err != nil {
		// return a 400 if the request body doesn't decode to a StaticInventory
		return err, http.StatusBadRequest
	}

	// pass the decoded model to the dao to update the DB
	static, err = c.dao.CreateStatic(static)

	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	// set computed field on static inventory model for thumbnail image URL
	static.SetThumbnailURL()

	// return the created static so that client can find out the auto generated ids
	c.JSON(rw, http.StatusOK, static)
	return nil, http.StatusOK
}
