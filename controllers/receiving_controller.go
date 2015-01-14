// Copyright G2G Market Inc, 2015

package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"gopkg.in/unrolled/render.v1"
	"just-pikd-wms/daos"
	"just-pikd-wms/helpers"
	"just-pikd-wms/models"
	"net/http"
	"strconv"
)

// ReceivingController contains methods related to the WMS Receiving module
// Struct members are passed in when declaring routes, except BaseController which is embedded for its methods
// TODO: think about whether it even makes sense to separate controllers by module, since different modules
// will likely share methods. Does it make sense to think that one module might "own" a method though?
// Perhaps instead use things like InventoryController, LocationController, etc.
type ReceivingController struct {
	BaseController
	*render.Render
	*sqlx.DB
}

// Reset is a helper function for resetting test data in dev only
// It calls the helper to reload data from json files
func (c *ReceivingController) Reset(rw http.ResponseWriter, r *http.Request) (error, int) {
	helpers.Reset(c.DB)
	rw.Write([]byte("Data reset complete"))
	return nil, http.StatusOK
}

// GetSPO gets a StockingPurchaseOrder from the database based on its id
// TODO: we'll have other ways of searching so possibly rename to GetSPOByID?
func (c *ReceivingController) GetSPO(rw http.ResponseWriter, r *http.Request) (error, int) {
	// parse args - no need to check error because gorilla/mux would 404 on invalid params anyway
	spo_id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// initialize dao and use it to retrieve spo model
	dao := daos.StockingPurchaseOrderDAO{DB: c.DB}
	spo, err := dao.GetSPO(spo_id)

	// handle errors - no need to log a 404 for rows not found, but do log other db errors
	if err == sql.ErrNoRows {
		return err, http.StatusNotFound
	} else if err != nil {
		c.LogError(err.Error())
		return err, http.StatusInternalServerError
	}

	// render the model as JSON as response
	c.JSON(rw, http.StatusOK, spo)
	return nil, http.StatusOK
}

/*
//Create inbound inventory record
//TODO parse JSON body to make this actually work
func (c *ReceivingController) CreateInbound(rw http.ResponseWriter, r *http.Request) (error, int) {
	now := time.Now()
	inbound := models.InboundInventory{1,
		null.NewInt(1, true),
		"1",
		&now,
		&now,
		2,
		null.NewInt(0, false),
		"Ordered",
		null.NewString("", false),
		null.NewString("", false)}
	dao := daos.InventoryDAO{DB: c.DB}
	inbound, err := dao.Create_Inbound(inbound)
	if err != nil {
		c.LogError(err.Error())
		return err, http.StatusInternalServerError
	}

	c.JSON(rw, http.StatusOK, inbound)
	return nil, http.StatusOK
}
*/

// GetStatic retrieves a static inventory record by its id
func (c *ReceivingController) GetStatic(rw http.ResponseWriter, r *http.Request) (error, int) {
	// parse args
	static_id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// retrieve static inventory model from dao by id
	dao := daos.InventoryDAO{DB: c.DB}
	static, err := dao.GetStatic(static_id)

	// handle errors
	if err == sql.ErrNoRows {
		return err, http.StatusNotFound
	} else if err != nil {
		c.LogError(err.Error())
		return err, http.StatusInternalServerError
	}

	// set computed field on static inventory model for thumbnail image URL
	static.SetThumbnailURL()

	// render and return response
	c.JSON(rw, http.StatusOK, static)
	return nil, http.StatusOK
}

// GetStockingLocation retrieves a stocking location based on its id
// note, this is actually part of stocking not receiving
func (c *ReceivingController) GetStockingLocation(rw http.ResponseWriter, r *http.Request) (error, int) {
	stocking_location_id := mux.Vars(r)["id"]

	dao := daos.LocationDAO{DB: c.DB}
	location, err := dao.GetStockingLocation(stocking_location_id)

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
func (c *ReceivingController) GetReceivingLocations(rw http.ResponseWriter, r *http.Request) (error, int) {
	hp := r.FormValue("has_product")
	temperature_zone := r.FormValue("temperature_zone")
	var has_product bool

	// accept string values 1 or true for has_product as true, other values considered false
	if hp == "1" || hp == "true" {
		has_product = true
	}

	// retrieve slice of locations from dao based on params
	dao := daos.LocationDAO{DB: c.DB}
	locations, err := dao.GetReceivingLocations(temperature_zone, has_product)

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
func (c *ReceivingController) UpdateReceivingLocation(rw http.ResponseWriter, r *http.Request) (error, int) {
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
	dao := daos.LocationDAO{DB: c.DB}
	err = dao.UpdateReceivingLocation(location)

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

// GetShipments retrieves an array of shipments based on passed in filters, or all shipments if no filters
func (c *ReceivingController) GetShipments(rw http.ResponseWriter, r *http.Request) (error, int) {
	// parse form values if passed in
	shipment_id := r.FormValue("shipment_id")
	spo := r.FormValue("stocking_purchase_order_id")
	var stocking_purchase_order_id int
	if len(spo) > 0 {
		var err error
		stocking_purchase_order_id, err = strconv.Atoi(spo)
		if err != nil {
			return err, http.StatusBadRequest
		}
	}

	// retrieve slice of shipments from dao based on params
	dao := daos.SupplierDAO{DB: c.DB}
	shipments, err := dao.GetShipments(shipment_id, stocking_purchase_order_id)

	if err == sql.ErrNoRows {
		return err, http.StatusNotFound
	} else if err != nil {
		c.LogError(err.Error())
		return err, http.StatusInternalServerError
	}

	c.JSON(rw, http.StatusOK, shipments)
	return nil, http.StatusOK
}

// UpdateShipment updates a supplier shipment based on a passed in model
// and is used to set the arrival date field when scanning in a received shipment
// based on supplier_shipment_id
func (c *ReceivingController) UpdateShipment(rw http.ResponseWriter, r *http.Request) (error, int) {
	shipment_id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// parse request body for a JSON receiving shipment model
	decoder := json.NewDecoder(r.Body)
	var shipment models.SupplierShipment
	err := decoder.Decode(&shipment)

	if err != nil {
		// return a 400 if the request body doesn't decode to a SupplierShipment
		return err, http.StatusBadRequest
	}

	if shipment_id != shipment.Id {
		return errors.New("Identifier does not match request body for supplier_shipment_id"), http.StatusBadRequest
	}

	// pass the decoded model to the dao to update the DB
	dao := daos.SupplierDAO{DB: c.DB}
	err = dao.UpdateShipment(shipment)

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
