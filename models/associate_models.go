// Copyright G2G Market Inc, 2015

package models

import (
	"gopkg.in/guregu/null.v2"
	"time"
)

type Associate struct {
	Id        int         `json:"as_id"`
	FirstName string      `json:"as_first_name"`
	LastName  string      `json:"as_last_name"`
	LoginPin  null.String `json:"as_login_pin"`
}

type AssociateStation struct {
	Id          int        `json:"ast_id"`
	AsId        int        `json:"ast_as_id"`
	StationType string     `json:"ast_station_type"`
	StartTime   *time.Time `json:"ast_start_time"`
	EndTime     *time.Time `json:"ast_end_time"`
}
