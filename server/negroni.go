// Copyright G2G Market Inc, 2015

package server

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
	// create new negroni middleware handler
	n := negroni.New()

	// gzip compression middleware
	n.Use(gzip.Gzip(gzip.DefaultCompression))

	// logging middleware
	n.Use(negroni.NewLogger())

	//static file serving middleware
	static := negroni.NewStatic(http.Dir(config.StaticDir))
	static.Prefix = "/" + config.StaticDir
	n.Use(static)

	// panic recovery middleware should be registered last
	// so that its deferred function runs before other middlewares
	n.Use(negroni.NewRecovery())

	//add the mux router as the handler for negroni
	n.UseHandler(router)

	return n
}
