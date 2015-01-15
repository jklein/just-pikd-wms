// Copyright G2G Market Inc, 2015

package controllers

import (
	"database/sql"
	//"encoding/json"
	//"errors"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"gopkg.in/unrolled/render.v1"
	"just-pikd-wms/daos"
	//"just-pikd-wms/models"
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
