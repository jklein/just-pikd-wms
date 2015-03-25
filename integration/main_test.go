// Copyright G2G Market Inc, 2015

package integration

import (
	"fmt"
	"just-pikd-wms/config"
	"just-pikd-wms/server"
	"net/http"
	"os"
	"testing"
	"time"
)

// This file runs integration tests, which use the entire service end-to-end including
// the real webserver and database

// TestMain handles test setup and then runs tests
func TestMain(m *testing.M) {
	//setup tasks - load config and create negroni instance
	c := config.New()
	n := server.New(c)

	//hack so that test data loads via relative path - start server in the main directory
	os.Chdir("../")

	//run the app in a background goroutine while testing continues
	go n.Run("localhost:3002")

	//sleep one second so that the server will be ready to listen for connections - without this the
	//below client intermittently fails with a connection refused error
	time.Sleep(time.Second)

	//start by invoking the handler to reset data, and make sure that succeeds
	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://localhost:3002/reset", nil)
	req.Header.Set("X-Auth-Token", server.STATIC_TOKEN)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error resetting data: %s", err.Error())
		os.Exit(1)
	} else if resp.StatusCode != 200 {
		fmt.Printf("Error resetting data: status %d", resp.StatusCode)
		os.Exit(1)
	}

	//once data has been reset, we're ready to run tests
	status := m.Run()
	//tear down
	os.Exit(status)
}
