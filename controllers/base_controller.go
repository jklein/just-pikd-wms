// Copyright G2G Market Inc, 2015

package controllers

import (
	"fmt"
	"net/http"
	"runtime"
)

// Action defines a standard function signature for us to use when creating
// controller actions. A controller action is basically just a method attached to
// a controller.
// We return err, int in that order because the int is only used if the error is not nil
// as the response code. Non-failure response codes should be set by the Action itself (like redirects)
type Action func(rw http.ResponseWriter, r *http.Request) (error, int)

// BaseController holds basic methods that will be used by all controllers, which embed this struct
type BaseController struct{}

// LogError contains basic error logging functionality available to controllers
// TODO add request context (TXID) type values
// TODO log to file and/or third party service (asynchronously with a goroutine)
func (c *BaseController) LogError(err string) {
	stack_trace := make([]byte, 1024)
	bytes_written := runtime.Stack(stack_trace, false)
	error_message := fmt.Sprintf("%s\n%s\n%d", err, stack_trace, bytes_written)
	fmt.Println(error_message)
}

// Action helps with error handling in a controller and wraps HandlerFuncs so that we can pass
// things like sql connections to them via their structs
func (c *BaseController) Action(a Action) http.HandlerFunc {
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
