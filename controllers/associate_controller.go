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

// associateController contains methods related to WMS associates
type associateController struct {
	baseController
	*render.Render
	*sqlx.DB
	dao daos.AssociateDAO
}

// NewLocationController acts as an initializer for associateController, setting
// its dao and returning the instance
func NewAssociateController(rend *render.Render, db *sqlx.DB) *associateController {
	c := new(associateController)
	c.Render = rend
	c.DB = db
	c.dao = daos.AssociateDAO{DB: db}
	return c
}

// Login handles a login event
func (c *associateController) Login(rw http.ResponseWriter, r *http.Request) (error, int) {
	type PinJSON struct {
		Pin string `json:"pin"`
	}

	var pj PinJSON
	err := jsonDecode(r.Body, &pj)

	if err != nil {
		return err, http.StatusBadRequest
	}

	pin := pj.Pin

	if len(pin) == 0 {
		return errors.New("Invalid request, missing pin parameter"), http.StatusBadRequest
	}

	//check PIN is valid -> find associate
	//  returns 404 if invalid pin
	associate, err := c.dao.GetAssociateByPin(pin)

	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	//TODO change Other to a constant by auto generating enums
	token, err := c.dao.CreateAssociateStation(associate.Id, "Other")

	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	// return associate id, name, auto generated auth token on response
	resp := map[string]string{"id": strconv.Itoa(associate.Id), "name": associate.FirstName, "token": token}
	c.JSON(rw, http.StatusOK, resp)
	return nil, 0
}

// Logout handles a logout event
func (c *associateController) Logout(rw http.ResponseWriter, r *http.Request) (error, int) {
	as_id, _ := strconv.Atoi(mux.Vars(r)["id"])
	token := r.Header.Get("X-Auth-Token")

	if len(token) == 0 {
		return errors.New("Invalid request, missing X-Auth-Token"), http.StatusBadRequest
	}

	// pass the decoded model to the dao to update the DB
	err := c.dao.EndAssociateSession(as_id, token)

	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	// no response needed on success
	return nil, http.StatusOK
}

// GetAssociate retrieves an associate by its id
func (c *associateController) GetAssociate(rw http.ResponseWriter, r *http.Request) (error, int) {
	as_id, _ := strconv.Atoi(mux.Vars(r)["id"])

	associate, err := c.dao.GetAssociate(as_id)

	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	c.JSON(rw, http.StatusOK, associate)
	return nil, 0
}

// UpdateAssociate updates an associate record based on a passed in JSON object
func (c *associateController) UpdateAssociate(rw http.ResponseWriter, r *http.Request) (error, int) {
	// extract identifier from url - while we don't use this, it helps follow REST principles to have it in the URI
	// and could later be used for something like varnish cache invalidation
	as_id, _ := strconv.Atoi(mux.Vars(r)["id"])
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return err, http.StatusBadRequest
	}

	var associate models.Associate
	err = json.Unmarshal(body, &associate)
	if err != nil {
		return err, http.StatusBadRequest
	} else if as_id != associate.Id {
		return errors.New("Identifier does not match request body for as_id"), http.StatusBadRequest
	}

	// also decode to a dict so that update statements can be handled
	dict, err := jsonToDict(body)
	if err != nil {
		return err, http.StatusBadRequest
	}

	// pass the decoded model to the dao to update the DB
	err = c.dao.UpdateAssociate(associate, dict)

	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	return nil, http.StatusNoContent
}

// CreateAssociate creates a new associate record based on a passed in JSON object
func (c *associateController) CreateAssociate(rw http.ResponseWriter, r *http.Request) (error, int) {
	var associate models.Associate
	err := jsonDecode(r.Body, &associate)

	if err != nil {
		// return a 400 if the request body doesn't decode to a ReceivingLocation
		return err, http.StatusBadRequest
	}

	// pass the decoded model to the dao to update the DB
	associate, err = c.dao.CreateAssociate(associate)

	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	c.JSON(rw, http.StatusCreated, associate)
	return nil, 0
}

func (c *associateController) SetAssociateStation(rw http.ResponseWriter, r *http.Request) (error, int) {
	as_id, _ := strconv.Atoi(mux.Vars(r)["id"])
	type StationJSON struct {
		Station string `json:"station"`
	}

	var sj StationJSON
	err := jsonDecode(r.Body, &sj)

	if err != nil {
		return err, http.StatusBadRequest
	}

	station := sj.Station

	if len(station) == 0 {
		return errors.New("Invalid request, missing station parameter"), http.StatusBadRequest
	}

	//TODO validate against list of enums
	token, err := c.dao.CreateAssociateStation(as_id, station)

	if err != nil {
		return err, c.sqlErrorToStatusCodeAndLog(err)
	}

	// return auto generated auth token on response
	c.JSON(rw, http.StatusCreated, map[string]string{"token": token})
	return nil, 0
}
