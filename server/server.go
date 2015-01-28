// Copyright G2G Market Inc, 2015

// Package server contains the webserver handlers
package server

import (
	"github.com/codegangsta/negroni"
	"just-pikd-wms/config"
)

// New initializes the server and returns a negroni instance ready to run
// TODO: consider wrapping this in our own struct and moving config loading in here
func New(config *config.Config) *negroni.Negroni {
	// create DB connection pool
	db := SetupDB(config.DbUser, config.DbPass, config.DbName)

	// set up routes
	router := MakeRouter(db, config)

	// create negroni middleware handler
	n := MakeNegroni(router, config)

	return n
}
