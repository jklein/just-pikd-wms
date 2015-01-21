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

// StockingPurchaseOrderController contains routes related to stocking purchase orders
type StockingPurchaseOrderController struct {
	BaseController
	*render.Render
	*sqlx.DB
}

// GetSPO gets a StockingPurchaseOrder from the database based on its id
// TODO: we'll have other ways of searching so possibly rename to GetSPOByID?
func (c *StockingPurchaseOrderController) GetSPO(rw http.ResponseWriter, r *http.Request) (error, int) {
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

// CreateSPO creates a new stocking purchase order along with its products
// based on a passed in JSON object
func (c *StockingPurchaseOrderController) CreateSPO(rw http.ResponseWriter, r *http.Request) (error, int) {
	// parse request body for a JSON receiving shipment model
	decoder := json.NewDecoder(r.Body)
	var spo models.StockingPurchaseOrder
	err := decoder.Decode(&spo)

	if err != nil {
		// return a 400 if the request body doesn't decode to a StockingPurchaseOrder
		return err, http.StatusBadRequest
	}

	// pass the decoded model to the dao to update the DB
	dao := daos.StockingPurchaseOrderDAO{DB: c.DB}
	spo, err = dao.CreateSPO(spo)

	if err != nil {
		c.LogError(err.Error())
		return err, http.StatusInternalServerError
	}

	// return the created spo so that client can find out the auto generated ids
	c.JSON(rw, http.StatusOK, spo)
	return nil, http.StatusOK
}

// CreateSPOProduct adds a new product to an existing SPO
// it returns a 404 if the SPO does not already exist
func (c *StockingPurchaseOrderController) CreateSPOProduct(rw http.ResponseWriter, r *http.Request) (error, int) {
	// parse request body for a JSON receiving shipment model
	spo_id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// first make sure the spo itself exists so that we can return a 404 otherwise
	dao := daos.StockingPurchaseOrderDAO{DB: c.DB}
	_, err := dao.GetSPO(spo_id)

	if err == sql.ErrNoRows {
		return err, http.StatusNotFound
	} else if err != nil {
		c.LogError(err.Error())
		return err, http.StatusInternalServerError
	}

	decoder := json.NewDecoder(r.Body)
	var spo_product models.StockingPurchaseOrderProduct
	err = decoder.Decode(&spo_product)

	// 400 if bad object, or id doesn't matched the identifier
	if err != nil || spo_product.SpoId != spo_id {
		return err, http.StatusBadRequest
	}

	// pass the decoded model to the dao to update the DB
	spo_product, err = dao.CreateSPOProduct(spo_product)

	if err != nil {
		c.LogError(err.Error())
		return err, http.StatusInternalServerError
	}

	// return the created object so client can extract id
	c.JSON(rw, http.StatusOK, spo_product)
	return nil, http.StatusOK
}

// UpdateSPO updates a stocking purchase order and/or its products based on passed in objects
func (c *StockingPurchaseOrderController) UpdateSPO(rw http.ResponseWriter, r *http.Request) (error, int) {
	spo_id, _ := strconv.Atoi(mux.Vars(r)["id"])
	// read request body to end and store so it can be unmarshaled into separate types later
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return err, http.StatusBadRequest
	}

	var spo models.StockingPurchaseOrder
	err = json.Unmarshal(body, &spo)

	if err != nil {
		return err, http.StatusBadRequest
	}

	if spo_id != spo.Id {
		return errors.New("Identifier does not match request body for stocking_purchase_order_id"), http.StatusBadRequest
	}

	// also decode to a dict so that update statements can be handled
	var dict map[string]interface{}
	err = json.Unmarshal(body, &dict)

	if err != nil {
		return err, http.StatusBadRequest
	}

	// pass the decoded model to the dao to update the DB
	dao := daos.StockingPurchaseOrderDAO{DB: c.DB}
	err = dao.UpdateSPO(spo, dict)

	if err == sql.ErrNoRows {
		return err, http.StatusNotFound
	} else if err != nil {
		c.LogError(err.Error())
		return err, http.StatusInternalServerError
	}

	// no response body needed for succesful update, just return 200
	return nil, http.StatusOK
}
