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
	"just-pikd-wms/models"
	"net/http"
	"strconv"
)

// SupplierController contains methods related to suppliers and the objects they own
type SupplierController struct {
	BaseController
	*render.Render
	*sqlx.DB
}

// GetShipments retrieves an array of shipments based on passed in filters, or all shipments if no filters
func (c *SupplierController) GetShipments(rw http.ResponseWriter, r *http.Request) (error, int) {
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

	// retrieve slice of shipments from dao based on params
	dao := daos.SupplierDAO{DB: c.DB}
	shipments, err := dao.GetShipments(shipment_code, spo_id)

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
func (c *SupplierController) UpdateShipment(rw http.ResponseWriter, r *http.Request) (error, int) {
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
