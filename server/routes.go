// Copyright G2G Market Inc, 2015

package server

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"gopkg.in/unrolled/render.v1"
	"just-pikd-wms/config"
	"just-pikd-wms/controllers"
)

// MakeRouter creates a gorilla/mux router and sets all routes to hit our controllers
func MakeRouter(db *sqlx.DB, config *config.Config) *mux.Router {
	router := mux.NewRouter()

	//create a render instance used to render JSON
	rend := render.New(render.Options{IndentJSON: true})

	//initialize controllers and then their routes
	sc := controllers.NewStockingPurchaseOrderController(rend, db)
	router.HandleFunc("/spos/{id:[0-9]+}", sc.Action(sc.GetSPO)).Methods("GET")
	router.HandleFunc("/spos", sc.Action(sc.GetSPOs)).Methods("GET")
	router.HandleFunc("/spos", sc.Action(sc.CreateSPO)).Methods("POST")
	router.HandleFunc("/spos/{id:[0-9]+}", sc.Action(sc.UpdateSPO)).Methods("PATCH")
	router.HandleFunc("/spos/{id:[0-9]+}/products", sc.Action(sc.CreateSPOProduct)).Methods("POST")
	router.HandleFunc("/spos/{id:[0-9]+}/products/{product_id:[0-9]+}", sc.Action(sc.GetSPOProduct)).Methods("GET")

	ic := controllers.NewInventoryController(rend, db)
	//TODO: remove /static/ from the route because there is only one kind of inventory now?
	router.HandleFunc("/inventory/static/{id:[0-9]+}", ic.Action(ic.GetStatic)).Methods("GET")
	router.HandleFunc("/inventory/static", ic.Action(ic.CreateStatic)).Methods("POST")
	router.HandleFunc("/inventory/static/{id:[0-9]+}", ic.Action(ic.UpdateStatic)).Methods("PATCH")
	//inventory by sku
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
	router.HandleFunc("/locations/containers/{id:[0-9-]+}", lc.Action(lc.GetPickContainer)).Methods("GET")
	router.HandleFunc("/locations/containers/{id:[0-9-]+}", lc.Action(lc.UpdatePickContainer)).Methods("PATCH")
	router.HandleFunc("/locations/containers", lc.Action(lc.CreatePickContainer)).Methods("POST")
	router.HandleFunc("/locations/container_locations", lc.Action(lc.CreatePickContainerLocation)).Methods("POST")
	router.HandleFunc("/locations/container_locations/{id:[0-9-]+}", lc.Action(lc.GetPickContainerLocation)).Methods("GET")
	router.HandleFunc("/locations/container_locations/{id:[0-9-]+}", lc.Action(lc.UpdatePickContainerLocation)).Methods("PATCH")
	router.HandleFunc("/locations/pickup", lc.Action(lc.CreatePickupLocation)).Methods("POST")
	router.HandleFunc("/locations/pickup/{id:[0-9]+}", lc.Action(lc.GetPickupLocation)).Methods("GET")
	router.HandleFunc("/locations/pickup/{id:[0-9]+}", lc.Action(lc.UpdatePickupLocation)).Methods("PATCH")
	//router.HandleFunc("/locations/kiosks", lc.Action(lc.CreateKiosk)).Methods("POST")
	//router.HandleFunc("/locations/kiosks/{id:[0-9]+}", lc.Action(lc.GetKiosk)).Methods("GET")
	//router.HandleFunc("/locations/kiosks/{id:[0-9]+}", lc.Action(lc.UpdateKiosk)).Methods("PATCH")

	uc := controllers.NewSupplierController(rend, db)
	router.HandleFunc("/suppliers/shipments", uc.Action(uc.GetShipments)).Methods("GET")
	router.HandleFunc("/suppliers/shipments/{id:[0-9]+}", uc.Action(uc.GetShipment)).Methods("GET")
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
	ac := controllers.NewAssociateController(rend, db)
	router.HandleFunc("/associates/{id:[0-9]+}/logout", ac.Action(ac.Logout)).Methods("POST")
	router.HandleFunc("/associates/{id:[0-9]+}", ac.Action(ac.GetAssociate)).Methods("GET")
	router.HandleFunc("/associates/{id:[0-9]+}", ac.Action(ac.UpdateAssociate)).Methods("PATCH")
	router.HandleFunc("/associates", ac.Action(ac.CreateAssociate)).Methods("POST")
	// /associates/{id:[0-9]+}/??? endpoint to reassign
	// /associates
	// /associate_stations - need something to track when an associate's station changes

	mainRouter := mux.NewRouter().StrictSlash(true)
	negroniAuth := negroni.New()
	negroniAuth.Use(NewAuthMiddleware(rend, db))
	negroniAuth.UseHandler(router)
	mainRouter.HandleFunc("/associates/login", ac.Action(ac.Login)).Methods("POST") // noAuth endpoint: for login only
	mainRouter.PathPrefix("/").Handler(negroniAuth)                                 //all other endpoints require auth

	//if you log in ->
	//clear any existing active logins
	//all requests require an auth token and login
	//expiring existing active logins would mean future web requests would fail
	//but what about resuming work in progress?

	return mainRouter
}
