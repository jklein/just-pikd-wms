// Copyright G2G Market Inc, 2015

package controllers

import (
	"database/sql"
	"errors"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"gopkg.in/unrolled/render.v1"
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
func (c *supplierController) UpdateShipment(rw http.ResponseWriter, r *http.Request) (error, int) {
	shipment_id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var shipment models.SupplierShipment
	err := jsonDecode(r.Body, &shipment)

	if err != nil {
		return err, http.StatusBadRequest
	}

	if shipment_id != shipment.Id {
		return errors.New("Identifier does not match request body for shi_id"), http.StatusBadRequest
	}

	err = c.dao.UpdateShipment(shipment)

	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	return nil, http.StatusOK
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
		c.LogError(err.Error())
		return err, http.StatusInternalServerError
	}

	// return the created shipment so that client can find out the auto generated ids
	c.JSON(rw, http.StatusOK, shipment)
	return nil, http.StatusOK
}
