//forked from github.com/ivpusic/httpcheck
//modified to not require a handler, and to allow specifying request body
package helpers

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

type (
	Checker struct {
		t        *testing.T
		addr     string
		request  *http.Request
		response *http.Response
		Body     []byte
		prefix   string
	}

	Callback func(*http.Response)
)

func H(t *testing.T) *Checker {
	return &Checker{
		t:      t,
		addr:   ":3002",
		prefix: "http://localhost:3002",
	}
}

// make request /////////////////////////////////////////////////

// If you want to provide you custom http.Request instance, you can do it using this method
// In this case internal http.Request instance won't be created, and passed instane will be used
// for making request
func (c *Checker) TestRequest(request *http.Request) *Checker {
	assert.NotNil(c.t, request, "Request nil")

	c.request = request
	return c
}

// Prepare for testing some part of code which lives on provided path and method.
func (c *Checker) Test(method, path string, body string) *Checker {
	method = strings.ToUpper(method)
	request, err := http.NewRequest(method, c.prefix+path, strings.NewReader(body))

	assert.Nil(c.t, err, "Failed to make new request")
	request.Header.Set("Content-Type", "application/json")
	request.Body.Close()

	c.request = request
	return c
}

// Final URL for request will be prefix+path.
// Prefix can be something like "http://localhost:3000", and path can be "/some/path" for example.
// Path is provided by user using "Test" method.
// Library will try to figure out URL prefix automatically for you.
// But in case that for your case is not the best, you can set prefix manually
func (c *Checker) SetPrefix(prefix string) *Checker {
	c.prefix = prefix
	return c
}

// headers ///////////////////////////////////////////////////////

// Will put header on request
func (c *Checker) WithHeader(key, value string) *Checker {
	c.request.Header.Set(key, value)
	return c
}

// Will check if response contains header on provided key with provided value
func (c *Checker) HasHeader(key, expectedValue string) *Checker {
	value := c.response.Header.Get(key)
	assert.Exactly(c.t, expectedValue, value)

	return c
}

// cookies ///////////////////////////////////////////////////////

// Will put cookie on request
func (c *Checker) HasCookie(key, expectedValue string) *Checker {
	found := false
	for _, cookie := range c.response.Cookies() {
		if cookie.Name == key && cookie.Value == expectedValue {
			found = true
			break
		}
	}
	assert.True(c.t, found)

	return c
}

// Will ckeck if response contains cookie with provided key and value
func (c *Checker) WithCookie(key, value string) *Checker {
	c.request.AddCookie(&http.Cookie{
		Name:  key,
		Value: value,
	})

	return c
}

// status ////////////////////////////////////////////////////////

// Will ckeck if response status is equal to provided
func (c *Checker) HasStatus(status int) *Checker {
	assert.Exactly(c.t, status, c.response.StatusCode)
	return c
}

// json body /////////////////////////////////////////////////////

// Will check if body contains json with provided value
func (c *Checker) HasJson(value interface{}) *Checker {
	valueBytes, err := json.Marshal(value)
	assert.Nil(c.t, err)
	assert.Equal(c.t, string(valueBytes), string(c.Body))
	return c
}

// body //////////////////////////////////////////////////////////

// Will check if body is equal to provided []byte data
func (c *Checker) HasBody(body []byte) *Checker {
	assert.Equal(c.t, body, c.Body)
	return c
}

// Will check if body contains string representation of provided value
func (c *Checker) BodyContains(str string) *Checker {
	assert.Contains(c.t, string(c.Body), str)
	return c
}

// Will make reqeust to built request object.
// After request is made, it will save response object for future assertions
// Responsibility of this method is also to start and stop HTTP server
func (c *Checker) Check() *Checker {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	response, err := client.Do(c.request)
	assert.Nil(c.t, err, "Failed while making new request.", err)

	// save response for assertion checks
	c.response = response
	body, err := ioutil.ReadAll(c.response.Body)
	assert.Nil(c.t, err)
	c.Body = body

	return c
}

// Will call provided callback function with current response
func (c *Checker) Cb(cb Callback) {
	cb(c.response)
}
