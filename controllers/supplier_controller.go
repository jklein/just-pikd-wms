// Copyright G2G Market Inc, 2015

package controllers

import (
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

// supplierController contains methods related to suppliers and the objects they own
// it's not exported so that it can't be instantiated without initializing the dao -
// instead it must be initialized by calling NewSupplierController
type supplierController struct {
	baseController
	*render.Render
	*sqlx.DB
	dao daos.SupplierDAO
}

// NewSupplierController acts as an initializer for supplierController, setting
// its dao and returning the instance
func NewSupplierController(rend *render.Render, db *sqlx.DB) *supplierController {
	c := new(supplierController)
	c.Render = rend
	c.DB = db
	c.dao = daos.SupplierDAO{DB: db}
	return c
}

// GetShipments retrieves an array of shipments based on passed in filters, or all shipments if no filters
func (c *supplierController) GetShipments(rw http.ResponseWriter, r *http.Request) (error, int) {
	// parse form values if passed in
	shipment_code := r.FormValue("shipment_code")
	spo := r.FormValue("spo_id")
	var spo_id int
	if len(spo) > 0 {
		var err error
		spo_id, err = strconv.Atoi(spo)
		if err != nil {
			return err, http.StatusBadRequest
		}
	}

	shipments, err := c.dao.GetShipments(shipment_code, spo_id)

	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	c.JSON(rw, http.StatusOK, shipments)
	return nil, 0
}

// GetShipment retrieves a single shipment record by its id
func (c *supplierController) GetShipment(rw http.ResponseWriter, r *http.Request) (error, int) {
	shi_id, _ := strconv.Atoi(mux.Vars(r)["id"])

	shipment, err := c.dao.GetShipment(shi_id)
	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	c.JSON(rw, http.StatusOK, shipment)
	return nil, 0
}

// CreateShipment creates a new supplier shipment record based on a passed in JSON object
func (c *supplierController) CreateShipment(rw http.ResponseWriter, r *http.Request) (error, int) {
	var shipment models.SupplierShipment
	err := jsonDecode(r.Body, &shipment)

	if err != nil {
		return err, http.StatusBadRequest
	}

	shipment, err = c.dao.CreateShipment(shipment)

	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	// return the created shipment so that client can find out the auto generated ids
	c.JSON(rw, http.StatusCreated, shipment)
	return nil, 0
}

// UpdateShipment updates a supplier shipment based on a passed in model
// and is used to set the arrival date field when scanning in a received shipment
// based on supplier_shipment_id
func (c *supplierController) UpdateShipment(rw http.ResponseWriter, r *http.Request) (error, int) {
	shipment_id, _ := strconv.Atoi(mux.Vars(r)["id"])
	// read request body to end and store so it can be unmarshaled into separate types later
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return err, http.StatusBadRequest
	}

	var shipment models.SupplierShipment
	err = json.Unmarshal(body, &shipment)

	if err != nil {
		return err, http.StatusBadRequest
	} else if shipment_id != shipment.Id {
		return errors.New("Identifier does not match request body for shi_id"), http.StatusBadRequest
	}

	// also decode to a dict so that update statements can be handled
	dict, err := jsonToDict(body)
	if err != nil {
		return err, http.StatusBadRequest
	}

	err = c.dao.UpdateShipment(shipment, dict)

	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	// no response body needed for succesful update, just return 200
	return nil, http.StatusNoContent
}
