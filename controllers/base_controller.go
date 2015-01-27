// Copyright G2G Market Inc, 2015

package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
)

// Action defines a standard function signature for us to use when creating
// controller actions. A controller action is basically just a method attached to
// a controller.
// We return err, int in that order because the int is only used if the error is not nil
// as the response code. Non-failure response codes should be set by the Action itself (like redirects)
type Action func(rw http.ResponseWriter, r *http.Request) (error, int)

// baseController holds basic methods that will be used by all controllers, which embed this struct
type baseController struct{}

// Controller is the interface that all controllers must satisfy
type Controller interface {
	Action(a Action) http.HandlerFunc
}

// LogError contains basic error logging functionality available to controllers
// TODO add request context (TXID) type values
// TODO log to file and/or third party service (asynchronously with a goroutine)
// TODO add ability to skip x stack frames
func (c *baseController) LogError(err string) {
	stack_trace := make([]byte, 1024)
	bytes_written := runtime.Stack(stack_trace, false)
	error_message := fmt.Sprintf("%s\n%s\n%d", err, stack_trace, bytes_written)
	fmt.Println(error_message)
}

// sqlErrorToStatusCode helps to categorize sql errors, converting sql.ErrNoRows to a 404 status code
// and other errors to a 500. It also logs 500s but does not log 404s
func (c *baseController) sqlErrorToStatusCodeAndLog(err error) int {
	if err == sql.ErrNoRows {
		// return 404 if the row was not found
		return http.StatusNotFound
	} else if err != nil {
		//TODO: skip 1 or 2 stack frames when logging this error
		c.LogError(err.Error())
		return http.StatusInternalServerError
	}
	// if we got here, a nil error was passed in
	return http.StatusOK
}

// Action helps with error handling in a controller and wraps HandlerFuncs so that we can pass
// things like sql connections to them via their structs
func (c *baseController) Action(a Action) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if err, status_code := a(rw, r); err != nil {
			//extract error message
			error_message := err.Error()

			//default to status code 500 on error
			if status_code == 0 {
				status_code = http.StatusInternalServerError
			}

			//overwrite the default error, which is often "sql: no rows in result set" in case of 404
			if status_code == http.StatusNotFound {
				error_message = "404 not found"
			}

			//return error message and code
			http.Error(rw, error_message, status_code)
		}
	})
}

// jsonDecode is a simple helper function to decode a response body into a struct
// to save a couple of lines of repetitive code in controllers
func jsonDecode(reader io.ReadCloser, v interface{}) error {
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(v)
	return err
}

// jsonToDict is a helper function for converting a json doc into a dictionary
func jsonToDict(input []byte) (map[string]interface{}, error) {
	var dict map[string]interface{}
	err := json.Unmarshal(input, &dict)
	return dict, err
}
