// Copyright G2G Market Inc, 2015

package server

import (
	"github.com/jmoiron/sqlx"
	"gopkg.in/unrolled/render.v1"
	"net/http"
)

type auth struct {
	*render.Render
	*sqlx.DB
}

//static token for dev/website API requests
const STATIC_TOKEN = "c268062d-7b99-46ba-8663-bde5870672aa"

func NewAuthMiddleware(rend *render.Render, db *sqlx.DB) *auth {
	m := new(auth)
	m.Render = rend
	m.DB = db
	return m
}

func (self *auth) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	//TODO come up with a separate place to store tokens for website API clients (generated during provisioning?)
	//and a header to specify which to use
	token := r.Header.Get("X-Auth-Token")

	if token == "" {
		//Authentication Information not included in request
		self.JSON(rw, http.StatusUnauthorized, "Missing header X-Auth-Token")
		return
	}

	if token == STATIC_TOKEN {
		//auth success, proceed
		next(rw, r)
	} else {

		var ast_id int
		err := self.DB.Get(&ast_id, "SELECT ast_id FROM associate_stations WHERE ast_api_token=$1 AND ast_end_time IS NULL", token)

		if err != nil || ast_id == 0 {
			self.JSON(rw, http.StatusUnauthorized, "Invalid X-Auth-Token or session expired")
			return
		}
		//auth success, proceed
		next(rw, r)
	}
}
