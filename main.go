// Copyright G2G Market Inc, 2015

package main

import (
	"just-pikd-wms/config"
	"just-pikd-wms/server"
)

// main initializes the server and then runs the app
func main() {
	// initialize and load config
	c := config.New()
	n := server.New(c)

	// run negroni to listen and serve the app
	n.Run(c.Host + ":" + c.Port)
}
