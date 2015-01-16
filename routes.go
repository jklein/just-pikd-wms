// Copyright G2G Market Inc, 2015

package main

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"gopkg.in/unrolled/render.v1"
	"just-pikd-wms/config"
	"just-pikd-wms/controllers"
)

func MakeRouter(db *sqlx.DB, config *config.Config) *mux.Router {
	//create gorilla/mux router
	router := mux.NewRouter()

	//create a render instance used to render JSON
	rend := render.New(render.Options{IndentJSON: true})

	//TODO: think about instantiating daos here and passing them to controllers so they can be mocked (DI)

	//initialize controllers and then their routes
	sc := &controllers.StockingPurchaseOrderController{Render: rend, DB: db}
	router.HandleFunc("/spos/{id:[0-9]+}", sc.Action(sc.GetSPO)).Methods("GET")
	//router.HandleFunc("/spos", sc.Action(sc.CreateSPO)).Methods("POST") NYI - purchasing
	//router.HandleFunc("/spos/{id:[0-9]+}", sc.Action(sc.UpdateSPO)).Methods("PUT") NYI - purchasing, stocking. this updates spo and its products.
	//TODO do we need /spos/{id}/products and /spos/{id}/products/{id} endpoints?

	ic := &controllers.InventoryController{Render: rend, DB: db}
	router.HandleFunc("/inventory/static/{id:[0-9]+}", ic.Action(ic.GetStatic)).Methods("GET")
	//router.HandleFunc("/inventory/static", ic.Action(ic.CreateStatic)).Methods("POST") - stocking
	//router.HandleFunc("/inventory/static", ic.Action(ic.UpdateStatic)).Methods("PUT") NYI - ?
	//router.HandleFunc("/inventory/outbound", c.Action(c.GetInbound)).Methods("GET") NYI - ?
	//router.HandleFunc("/inventory/outbound", c.Action(c.CreateInbound)).Methods("POST") NYI - ?
	//router.HandleFunc("/inventory/outbound", c.Action(c.UpdateInbound)).Methods("PUT") NYI - ?
	// /inventory/errors
	if config.IsDev {
		//register dangerous route to reset data in dev only
		router.HandleFunc("/reset", ic.Action(ic.Reset)).Methods("POST")
	}

	lc := &controllers.LocationController{Render: rend, DB: db}
	router.HandleFunc("/locations/stocking/{id:[0-9-]+}", lc.Action(lc.GetStockingLocation)).Methods("GET")
	//router.HandleFunc("/locations/stocking", lc.Action(lc.CreateStockingLocation)).Methods("POST") NYI - store setup
	//router.HandleFunc("/locations/stocking/{id:[0-9-]+}", lc.Action(lc.UpdateStockingLocation)).Methods("PUT") NYI - store setup
	router.HandleFunc("/locations/receiving", lc.Action(lc.GetReceivingLocations)).Methods("GET").Queries("temperature_zone", "{temperature_zone:[a-z]+}")
	//router.HandleFunc("/locations/receiving", lc.Action(lc.CreateReceivingLocations)).Methods("POST") NYI - store setup
	router.HandleFunc("/locations/receiving/{id:[0-9-]+}", lc.Action(lc.UpdateReceivingLocation)).Methods("PATCH")
	//router.HandleFunc("/locations/containers/{id:[0-9-]+}", lc.Action(lc.GetPickContainer)).Methods("GET") NYI - also need put
	//router.HandleFunc("/locations/containers/{id:[0-9-]+}", lc.Action(lc.UpdatePickContainer)).Methods("PATCH") NYI - set location in picking
	//router.HandleFunc("/locations/containers", lc.Action(lc.CreatePickContainer)).Methods("POST") NYI - store setup
	//router.HandleFunc("/locations/container_locations", lc.Action(lc.CreatePickContainerLocation)).Methods("POST") NYI - store setup
	//router.HandleFunc("/locations/container_locations/{id:[0-9-]+}", lc.Action(lc.GetPickContainerLocation)).Methods("GET") NYI - picking
	//pickup locations: GET, POST, PATCH?
	//kiosks: GET, POST?

	uc := &controllers.SupplierController{Render: rend, DB: db}
	router.HandleFunc("/suppliers/shipments", uc.Action(uc.GetShipments)).Methods("GET")
	//router.HandleFunc("/suppliers/shipments/{id:[0-9]+}", uc.Action(uc.GetShipment)).Methods("GET") NYI - ?
	//router.HandleFunc("/suppliers/shipments", uc.Action(uc.CreateShipment)).Methods("POST") NYI - purchasing
	router.HandleFunc("/suppliers/shipments/{id:[0-9]+}", uc.Action(uc.UpdateShipment)).Methods("PATCH")

	//orders controller
	//orders - GET, POST, PUT. updates the order and its products.
	//TODO do we need /orders/{id}/products and /orders/{id}/products/{id} endpoints?

	//tasks controller
	// /tasks/pick
	// /tasks/pickup

	//associates controller
	// /associates
	// /associates/stations
	return router
}
