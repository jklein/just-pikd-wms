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
	"just-pikd-wms/helpers"
	//"just-pikd-wms/models"
	"net/http"
	"strconv"
)

// ReceivingController contains handlers for /inventory routes
type InventoryController struct {
	BaseController
	*render.Render
	*sqlx.DB
}

// Reset is a helper function for resetting test data in dev only
// It calls the helper to reload data from json files
func (c *InventoryController) Reset(rw http.ResponseWriter, r *http.Request) (error, int) {
	helpers.Reset(c.DB)
	rw.Write([]byte("Data reset complete"))
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
func (c *InventoryController) GetStatic(rw http.ResponseWriter, r *http.Request) (error, int) {
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
