// Copyright G2G Market Inc, 2015

// Package config contains a type for configuration values used throughout the app
package config

import (
	"os"
)

// Config holds global configuration values
type Config struct {
	Port      string
	Host      string
	IsDev     bool
	DbUser    string
	DbPass    string
	DbName    string
	StaticDir string
}

// New is a convenience method to instantiate a config and load it
// Panics if config fails to load, because that means the app can't start anyway
func New() *Config {
	c := new(Config)
	err := c.Load()
	if err != nil {
		panic("Unable to load config: " + err.Error())
	}
	return c
}

// Load loads and parses configuration values from environment variables
func (c *Config) Load() error {
	c.IsDev = isDev()

	// parse port to listen on from ENV variables
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}
	c.Port = port

	// parse HOST to listen on from ENV variables - negroni will default to localhost if this is empty
	host := os.Getenv("HOST")
	if len(host) == 0 {
		host = "localhost"
	}
	c.Host = host

	// TODO load these from a configuration file (with encryption for the password)
	// possibly with chef data bags?
	// parametrize config file name(s) too
	c.DbUser = "postgres"
	c.DbPass = "justpikd"
	c.DbName = "wms_1"

	c.StaticDir = "public"

	return nil
}

// IsDev tells us whether we're in dev or not
// currently based on environment username being set to vagrant or not,
// which may need to change at some point
func isDev() bool {
	user := os.Getenv("USER")
	if user == "vagrant" {
		return true
	}
	return false
}
