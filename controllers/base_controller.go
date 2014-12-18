package controllers

import (
	"fmt"
	"net/http"
	"runtime"
)

// Action defines a standard function signature for us to use when creating
// controller actions. A controller action is basically just a method attached to
// a controller.
type Action func(rw http.ResponseWriter, r *http.Request) (error, int)

// This is our Base Controller
type BaseController struct{}

// Basic error logging function
// TODO add request context (TXID) type values
// TODO log to file and/or third party service (asynchronously with a goroutine)
func (c *BaseController) LogError(err string) {
	stack_trace := make([]byte, 1024)
	bytes_written := runtime.Stack(stack_trace, false)
	error_message := fmt.Sprintf("%s\n%s\n%d", err, stack_trace, bytes_written)
	fmt.Println(error_message)
}

// The action function helps with error handling in a controller
func (c *BaseController) Action(a Action) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if err, status_code := a(rw, r); err != nil {
			if status_code == 0 {
				status_code = http.StatusInternalServerError
			}
			http.Error(rw, err.Error(), status_code)
		}
	})
}
