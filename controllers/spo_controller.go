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

// stockingPurchaseOrderController contains routes related to stocking purchase orders
type stockingPurchaseOrderController struct {
	baseController
	*render.Render
	*sqlx.DB
	dao daos.StockingPurchaseOrderDAO
}

// NewStockingPurchaseOrderController acts as an initializer for stockingPurchaseOrderController, setting
// its dao and returning the instance
func NewStockingPurchaseOrderController(rend *render.Render, db *sqlx.DB) *stockingPurchaseOrderController {
	c := new(stockingPurchaseOrderController)
	c.Render = rend
	c.DB = db
	c.dao = daos.StockingPurchaseOrderDAO{DB: db}
	return c
}

// GetSPO gets a StockingPurchaseOrder from the database based on its id
func (c *stockingPurchaseOrderController) GetSPO(rw http.ResponseWriter, r *http.Request) (error, int) {
	// parse args - no need to check error because gorilla/mux would 404 on invalid params anyway
	spo_id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// initialize dao and use it to retrieve spo model
	spo, err := c.dao.GetSPO(spo_id)

	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	// render the model as JSON as response
	c.JSON(rw, http.StatusOK, spo)
	return nil, 0
}

// GetSPOs searches for an array of spos based on passed in filters
func (c *stockingPurchaseOrderController) GetSPOs(rw http.ResponseWriter, r *http.Request) (error, int) {
	shipment_code := r.FormValue("shipment_code")
	su_id := r.FormValue("supplier_id")
	//TODO add shi_id as well or instead?
	//also consider adding the ability to filter by date
	var supplier_id int
	if len(su_id) > 0 {
		var err error
		supplier_id, err = strconv.Atoi(su_id)
		if err != nil {
			return err, http.StatusBadRequest
		}
	}

	// retrieve slice of locations from dao based on params
	spos, err := c.dao.GetSPOs(supplier_id, shipment_code)

	//no need to throw error if no rows are found as that is a normal case here
	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	c.JSON(rw, http.StatusOK, spos)
	return nil, 0
}

// CreateSPO creates a new stocking purchase order along with its products
// based on a passed in JSON object
func (c *stockingPurchaseOrderController) CreateSPO(rw http.ResponseWriter, r *http.Request) (error, int) {
	var spo models.StockingPurchaseOrder
	err := jsonDecode(r.Body, &spo)

	if err != nil {
		// return a 400 if the request body doesn't decode to a StockingPurchaseOrder
		return err, http.StatusBadRequest
	}

	spo, err = c.dao.CreateSPO(spo)

	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	// return the created spo so that client can find out the auto generated ids
	c.JSON(rw, http.StatusCreated, spo)
	return nil, 0
}

// GetSPO gets a StockingPurchaseOrder from the database based on its id
func (c *stockingPurchaseOrderController) GetSPOProduct(rw http.ResponseWriter, r *http.Request) (error, int) {
	// parse args - no need to check error because gorilla/mux would 404 on invalid params anyway
	spo_id, _ := strconv.Atoi(mux.Vars(r)["id"])
	spop_id, _ := strconv.Atoi(mux.Vars(r)["product_id"])

	// initialize dao and use it to retrieve spo model
	product, err := c.dao.GetSPOProduct(spo_id, spop_id)

	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	// render the model as JSON as response
	c.JSON(rw, http.StatusOK, product)
	return nil, 0
}

// CreateSPOProduct adds a new product to an existing SPO
// it returns a 404 if the SPO does not already exist
func (c *stockingPurchaseOrderController) CreateSPOProduct(rw http.ResponseWriter, r *http.Request) (error, int) {
	// parse request body for a JSON receiving shipment model
	spo_id, _ := strconv.Atoi(mux.Vars(r)["id"])

	// first make sure the spo itself exists so that we can return a 404 otherwise
	_, err := c.dao.GetSPO(spo_id)
	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	var spo_product models.StockingPurchaseOrderProduct
	err = jsonDecode(r.Body, &spo_product)
	if err != nil {
		return err, http.StatusBadRequest
	} else if spo_product.SpoId != spo_id {
		return errors.New("Identifier does not match request body for spop_spo_id"), http.StatusBadRequest
	}

	// pass the decoded model to the dao to update the DB
	spo_product, err = c.dao.CreateSPOProduct(spo_product)
	if err != nil {
		c.LogError(err.Error())
		return err, http.StatusInternalServerError
	}

	// return the created object so client can extract id
	c.JSON(rw, http.StatusCreated, spo_product)
	return nil, 0
}

// UpdateSPO updates a stocking purchase order and/or its products based on passed in objects
func (c *stockingPurchaseOrderController) UpdateSPO(rw http.ResponseWriter, r *http.Request) (error, int) {
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
	} else if spo_id != spo.Id {
		return errors.New("Identifier does not match request body for spo_id"), http.StatusBadRequest
	}

	// also decode to a dict so that update statements can be handled
	dict, err := jsonToDict(body)
	if err != nil {
		return err, http.StatusBadRequest
	}

	err = c.dao.UpdateSPO(spo, dict)

	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	// no response body needed for succesful update, just return 200
	return nil, http.StatusNoContent
}
