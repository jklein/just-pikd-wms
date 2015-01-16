// Copyright G2G Market Inc, 2015

package main

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/phyber/negroni-gzip/gzip"
	"just-pikd-wms/config"
	"net/http"
)

// MakeNegroni handles all negroni middleware configuration and returns a negroni instance
// ready to run
func MakeNegroni(router *mux.Router, config *config.Config) *negroni.Negroni {
	//create new negroni middleware handler
	n := negroni.New()

	//start with panic recovery middleware
	n.Use(negroni.NewRecovery())

	//logging middleware
	n.Use(negroni.NewLogger())

	//gzip compression middleware
	n.Use(gzip.Gzip(gzip.DefaultCompression))

	//static file serving middleware
	static := negroni.NewStatic(http.Dir(config.StaticDir))
	static.Prefix = "/" + config.StaticDir
	n.Use(static)

	//add the mux router as the handler for negroni
	n.UseHandler(router)

	return n
}
