package controllers

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"gopkg.in/unrolled/render.v1"
	"just-pikd-wms/daos"
	"just-pikd-wms/helpers"
	"just-pikd-wms/models"
	"net/http"
	"strconv"
)

type ReceivingController struct {
	BaseController
	*render.Render
	*sqlx.DB
}

func (c *ReceivingController) Reset(rw http.ResponseWriter, r *http.Request) (error, int) {
	helpers.Reset(c.DB)
	rw.Write([]byte("Data reset complete"))
	return nil, http.StatusOK
}

func (c *ReceivingController) GetSPO(rw http.ResponseWriter, r *http.Request) (error, int) {
	vars := mux.Vars(r)
	spo_id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err, http.StatusBadRequest
	}
	var spo models.StockingPurchaseOrder
	dao := daos.StockingPurchaseOrder_DAO{DB: c.DB}
	spo, err = dao.Get_SPO(spo_id)
	if err == sql.ErrNoRows {
		//no need to log 404?
		return err, http.StatusNotFound
	} else if err != nil {
		c.LogError(err.Error())
		return err, http.StatusInternalServerError
	}

	c.JSON(rw, http.StatusOK, spo)
	return nil, http.StatusOK
}

/*
//Create inbound inventory record
//TODO parse JSON body to make this real
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
	dao := daos.Inventory_DAO{DB: c.DB}
	inbound, err := dao.Create_Inbound(inbound)
	if err != nil {
		c.LogError(err.Error())
		return err, http.StatusInternalServerError
	}

	c.JSON(rw, http.StatusOK, inbound)
	return nil, http.StatusOK
}
*/

func (c *ReceivingController) GetStatic(rw http.ResponseWriter, r *http.Request) (error, int) {
	vars := mux.Vars(r)
	static_id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return err, http.StatusBadRequest
	}
	var static models.StaticInventory
	dao := daos.Inventory_DAO{DB: c.DB}
	static, err = dao.Get_Static(static_id)
	if err == sql.ErrNoRows {
		//no need to log 404?
		return err, http.StatusNotFound
	} else if err != nil {
		c.LogError(err.Error())
		return err, http.StatusInternalServerError
	}

	c.JSON(rw, http.StatusOK, static)
	return nil, http.StatusOK
}
