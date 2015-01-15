// Copyright G2G Market Inc, 2015

package main

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/phyber/negroni-gzip/gzip"
	"gopkg.in/unrolled/render.v1"
	"just-pikd-wms/controllers"
	"net/http"
	"os"
)

func main() {
	//initialize database
	//note - this pools connections automatically and opens them as needed when used (see documentation)
	//TODO figure out a way to securely store and retrieve password
	db, err := sqlx.Open("postgres", "user=postgres password='justpikd' dbname=wms_1 sslmode=disable")
	if err != nil {
		//No reason this should error here (even if the database is down or doesn't exist)
		panic(err)
	}

	//consider it to be in dev if user is vagrant - this may need to change at some point
	user := os.Getenv("USER")
	var is_dev bool
	if user == "vagrant" {
		is_dev = true
	}

	//create gorilla/mux router
	router := mux.NewRouter()

	//create a render instance used to render JSON
	rend := render.New(render.Options{IndentJSON: true})

	//initialize controllers and then their routes
	sc := &controllers.StockingPurchaseOrderController{Render: rend, DB: db}
	router.HandleFunc("/spos/{id:[0-9]+}", sc.Action(sc.GetSPO)).Methods("GET")

	ic := &controllers.InventoryController{Render: rend, DB: db}
	router.HandleFunc("/inventory/static/{id:[0-9]+}", ic.Action(ic.GetStatic)).Methods("GET")
	//router.HandleFunc("/inventory/inbound", c.Action(c.CreateInbound)).Methods("POST")
	if is_dev {
		//register dangerous route to reset data in dev only
		router.HandleFunc("/reset", ic.Action(ic.Reset)).Methods("POST")
	}

	lc := &controllers.LocationController{Render: rend, DB: db}
	router.HandleFunc("/locations/stocking/{id:[0-9-]+}", lc.Action(lc.GetStockingLocation)).Methods("GET")
	router.HandleFunc("/locations/receiving", lc.Action(lc.GetReceivingLocations)).Methods("GET").Queries("temperature_zone", "{temperature_zone:[a-z]+}")
	router.HandleFunc("/locations/receiving/{id:[0-9-]+}", lc.Action(lc.UpdateReceivingLocation)).Methods("PATCH")

	uc := &controllers.SupplierController{Render: rend, DB: db}
	router.HandleFunc("/suppliers/shipments", uc.Action(uc.GetShipments)).Methods("GET")
	router.HandleFunc("/suppliers/shipments/{id:[0-9]+}", uc.Action(uc.UpdateShipment)).Methods("PATCH")

	//parse port to listen on from ENV variables
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	//parse HOST to listen on from ENV variables - defaults to localhost
	host := os.Getenv("HOST")

	//create new negroni middleware handler
	n := negroni.New()

	//start with panic recovery middleware
	n.Use(negroni.NewRecovery())

	//logging middleware
	n.Use(negroni.NewLogger())

	//gzip compression middleware
	n.Use(gzip.Gzip(gzip.DefaultCompression))

	//static file serving middleware
	static := negroni.NewStatic(http.Dir("public"))
	static.Prefix = "/public"
	n.Use(static)

	//add the mux router as the handler for negroni
	n.UseHandler(router)
	//start the app
	n.Run(host + ":" + port)
}
