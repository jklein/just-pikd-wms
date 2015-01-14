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

	//initialize controllers
	rc := &controllers.ReceivingController{Render: rend, DB: db, Dev: is_dev}
	router.HandleFunc("/spo/{id:[0-9]+}", rc.Action(rc.GetSPO)).Methods("GET")
	//router.HandleFunc("/inventory/inbound", rc.Action(rc.CreateInbound)).Methods("POST")
	router.HandleFunc("/inventory/static/{id:[0-9]+}", rc.Action(rc.GetStatic)).Methods("GET")
	router.HandleFunc("/locations/stocking/{id:[0-9-]+}", rc.Action(rc.GetStockingLocation)).Methods("GET")
	router.HandleFunc("/locations/receiving", rc.Action(rc.GetReceivingLocations)).Methods("GET").Queries("temperature_zone", "{temperature_zone:[a-z]+}")
	router.HandleFunc("/locations/receiving", rc.Action(rc.PutReceivingLocation)).Methods("PUT")
	//DEV ONLY: route to reset all data
	router.HandleFunc("/reset", rc.Action(rc.Reset)).Methods("POST")

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
