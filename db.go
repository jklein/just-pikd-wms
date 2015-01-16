// Copyright G2G Market Inc, 2015

package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// SetupDB initializess the database connection
// note - this pools connections automatically and opens them as needed when used (see database/sql documentation)
func SetupDB(user string, pass string, dbname string) *sqlx.DB {
	conn_string := fmt.Sprintf("user=%s password='%s' dbname=%s sslmode=disable", user, pass, dbname)
	db, err := sqlx.Open("postgres", conn_string)
	if err != nil {
		//No reason this should error here (even if the database is down or doesn't exist)
		panic(err)
	}
	return db
}
