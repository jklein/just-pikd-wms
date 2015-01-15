// Copyright G2G Market Inc, 2015

package main

import (
	"just-pikd-wms/config"
)

// main initializes global objects and then runs the app
func main() {
	// initialize and load config
	config := &config.Config{}
	config.Load()

	// create DB connection pool
	db := SetupDB(config.DbUser, config.DbPass, config.DbName)

	// set up routes
	router := MakeRouter(db, config)

	// create negroni middleware handler
	n := MakeNegroni(router, config)

	// run negroni to listen and serve the app
	n.Run(config.Host + ":" + config.Port)
}
