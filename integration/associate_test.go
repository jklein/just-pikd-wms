// Copyright G2G Market Inc, 2015

package integration

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAssociate(t *testing.T) {
	H(t).Test("GET", "/associates/1", "").Check().HasStatus(200).BodyContains(`"as_first_name": "Scott"`)
}

func TestGetAssociateNotFound(t *testing.T) {
	H(t).Test("GET", "/associates/1000", "").Check().HasStatus(404)
}

func TestPostAssociate(t *testing.T) {
	body := `{"as_first_name": "Test", "as_last_name": "McTestington", "as_login_pin": "019192"}`
	H(t).Test("POST", "/associates", body).Check().HasStatus(201).BodyContains("as_id")
}

func TestPatchAssociate(t *testing.T) {
	body := `{"as_id": 1, "as_login_pin": "000000"}`
	H(t).Test("PATCH", "/associates/1", body).Check().HasStatus(204)
	H(t).Test("GET", "/associates/1", "").Check().HasStatus(200).BodyContains(`"as_login_pin": "000000"`)
}

func TestPatchAssociateIdentifierMismatch(t *testing.T) {
	body := `{"as_id": 2, "as_login_pin": "000000"}`
	H(t).Test("PATCH", "/associates/1", body).Check().HasStatus(400)
}

func TestPatchAssociateNonexistent(t *testing.T) {
	body := `{"as_id": 1000, "as_login_pin": "000200"}`
	H(t).Test("PATCH", "/associates/1000", body).Check().HasStatus(404)
}

func TestLoginAndLogout(t *testing.T) {
	//easier to test both in one test since logout requires the token that login produces
	type AuthResponse struct {
		Token string `json:"X-Auth-Token"`
	}
	body := `{"pin": "12346"}`
	respbody := H(t).Test("POST", "/associates/login", body).Check().HasStatus(200).BodyContains(`"X-Auth-Token"`).Body
	var resp AuthResponse
	err := json.Unmarshal(respbody, &resp)
	assert.Nil(t, err)

	H(t).Test("POST", "/associates/2/logout", "").WithHeader("X-Auth-Token", resp.Token).Check().HasStatus(200)
}

func TestInvalidLogin(t *testing.T) {
	body := `{"pin": "12347"}`
	H(t).Test("POST", "/associates/login", body).Check().HasStatus(404)
}
