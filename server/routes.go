// Copyright G2G Market Inc, 2015

package server

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"gopkg.in/unrolled/render.v1"
	"just-pikd-wms/config"
	"just-pikd-wms/controllers"
)

// MakeRouter creates a gorilla/mux router and sets all routes to hit our controllers
func MakeRouter(db *sqlx.DB, config *config.Config) *mux.Router {
	//create gorilla/mux router
	router := mux.NewRouter()

	//create a render instance used to render JSON
	rend := render.New(render.Options{IndentJSON: true})

	//TODO: think about instantiating daos here and passing them to controllers so they can be mocked (DI)

	//initialize controllers and then their routes
	sc := controllers.NewStockingPurchaseOrderController(rend, db)
	router.HandleFunc("/spos/{id:[0-9]+}", sc.Action(sc.GetSPO)).Methods("GET")
	router.HandleFunc("/spos", sc.Action(sc.GetSPOs)).Methods("GET")
	router.HandleFunc("/spos", sc.Action(sc.CreateSPO)).Methods("POST")
	router.HandleFunc("/spos/{id:[0-9]+}", sc.Action(sc.UpdateSPO)).Methods("PATCH")
	router.HandleFunc("/spos/{id:[0-9]+}/products", sc.Action(sc.CreateSPOProduct)).Methods("POST")
	//router.HandleFunc("/spos/{id:[0-9]+}/products/{product_id:[0-9]+}", sc.Action(sc.GetSPOProduct)).Methods("GET") NYI - not sure if we need

	ic := controllers.NewInventoryController(rend, db)
	//TODO: remove /static/ from the route because there is only one kind of inventory now?
	router.HandleFunc("/inventory/static/{id:[0-9]+}", ic.Action(ic.GetStatic)).Methods("GET")
	router.HandleFunc("/inventory/static", ic.Action(ic.CreateStatic)).Methods("POST")
	router.HandleFunc("/inventory/static/{id:[0-9]+}", ic.Action(ic.UpdateStatic)).Methods("PATCH")
	// /inventory/errors
	// /inventory/holds
	if config.IsDev {
		//register dangerous route to reset data in dev only
		router.HandleFunc("/reset", ic.Action(ic.Reset)).Methods("POST")
	}

	lc := controllers.NewLocationController(rend, db)
	router.HandleFunc("/locations/stocking/{id:[0-9-]+}", lc.Action(lc.GetStockingLocation)).Methods("GET")
	router.HandleFunc("/locations/stocking", lc.Action(lc.CreateStockingLocation)).Methods("POST")
	router.HandleFunc("/locations/stocking/{id:[0-9-]+}", lc.Action(lc.UpdateStockingLocation)).Methods("PATCH")
	router.HandleFunc("/locations/receiving", lc.Action(lc.GetReceivingLocations)).Methods("GET").Queries("temperature_zone", "{temperature_zone:[a-z]+}")
	router.HandleFunc("/locations/receiving/{id:[0-9-]+}", lc.Action(lc.GetReceivingLocation)).Methods("GET")
	router.HandleFunc("/locations/receiving", lc.Action(lc.CreateReceivingLocation)).Methods("POST")
	router.HandleFunc("/locations/receiving/{id:[0-9-]+}", lc.Action(lc.UpdateReceivingLocation)).Methods("PATCH")
	//router.HandleFunc("/locations/containers/{id:[0-9-]+}", lc.Action(lc.GetPickContainer)).Methods("GET") NYI - also need put
	//router.HandleFunc("/locations/containers/{id:[0-9-]+}", lc.Action(lc.UpdatePickContainer)).Methods("PATCH") NYI - set location in picking
	//router.HandleFunc("/locations/containers", lc.Action(lc.CreatePickContainer)).Methods("POST") NYI - store setup
	//router.HandleFunc("/locations/container_locations", lc.Action(lc.CreatePickContainerLocation)).Methods("POST") NYI - store setup
	//router.HandleFunc("/locations/container_locations/{id:[0-9-]+}", lc.Action(lc.GetPickContainerLocation)).Methods("GET") NYI - picking
	// /locations/pickup: GET, POST, PATCH?
	// /kiosks: GET, POST?

	uc := controllers.NewSupplierController(rend, db)
	router.HandleFunc("/suppliers/shipments", uc.Action(uc.GetShipments)).Methods("GET")
	//router.HandleFunc("/suppliers/shipments/{id:[0-9]+}", uc.Action(uc.GetShipment)).Methods("GET") NYI - ?
	router.HandleFunc("/suppliers/shipments", uc.Action(uc.CreateShipment)).Methods("POST")
	router.HandleFunc("/suppliers/shipments/{id:[0-9]+}", uc.Action(uc.UpdateShipment)).Methods("PATCH")
	// /suppliers/{id}

	//orders controller
	//orders (embedded object)
	//TODO do we need /orders/{id}/products and /orders/{id}/products/{id} endpoints?

	//tasks controller
	// /tasks/pick (embedded object)
	// /tasks/pickup (embedded object)

	//associates controller
	// /associates
	// /associate_stations
	return router
}
