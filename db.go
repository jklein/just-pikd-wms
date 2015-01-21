// Copyright G2G Market Inc, 2015

package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	_ "github.com/lib/pq"
	"strings"
)

// SetupDB initializes the database connection
// note - this pools connections automatically and opens them as needed when used (see database/sql documentation)
func SetupDB(user string, pass string, dbname string) *sqlx.DB {
	conn_string := fmt.Sprintf("user=%s password='%s' dbname=%s sslmode=disable", user, pass, dbname)
	db, err := sqlx.Open("postgres", conn_string)
	if err != nil {
		//No reason this should error here (even if the database is down or doesn't exist)
		panic(err)
	}
	db.Mapper = reflectx.NewMapperFunc("json", strings.ToLower)
	return db
}
