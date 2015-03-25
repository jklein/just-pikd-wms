// Copyright G2G Market Inc, 2015

package daos

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"just-pikd-wms/models"
)

// AssociateDAO is used for data access related to
// associates such as login, logout, creating/updating associates
type AssociateDAO struct {
	*sqlx.DB
}

// GetAssociate retrieves an associate record by its id
func (dao *AssociateDAO) GetAssociate(as_id int) (models.Associate, error) {
	var associate models.Associate

	err := dao.DB.Get(&associate,
		`SELECT as_id, as_first_name, as_last_name, as_login_pin
        FROM associates
        WHERE as_id = $1;`, as_id)
	return associate, err
}

// GetAssociateByPin retrieves an associate record by the PIN
func (dao *AssociateDAO) GetAssociateByPin(pin string) (models.Associate, error) {
	var associate models.Associate

	err := dao.DB.Get(&associate,
		`SELECT as_id, as_first_name, as_last_name, as_login_pin
        FROM associates
        WHERE as_login_pin = $1;`, pin)
	return associate, err
}

// UpdateAssociate updates an associate, updating only the passed-in fields
func (dao *AssociateDAO) UpdateAssociate(associate models.Associate, dict map[string]interface{}) error {
	stmt := buildPatchUpdate("associates", "as_id", dict)
	err := execCheckRows(dao.DB, stmt, associate)
	return err
}

// CreateAssociate creates a new associate record based on the passed in model
func (dao *AssociateDAO) CreateAssociate(associate models.Associate) (models.Associate, error) {
	rows, err := dao.DB.NamedQuery(
		`INSERT INTO associates (as_first_name, as_last_name, as_login_pin)
        VALUES (:as_first_name, :as_last_name, :as_login_pin)
        RETURNING as_id`,
		associate)
	if err != nil {
		return associate, err
	}
	id, err := extractLastInsertId(rows)
	associate.Id = id
	return associate, err
}

// GetAssociateStation retrieves an associate's current station
func (dao *AssociateDAO) GetAssociateStation(as_id int) (models.AssociateStation, error) {
	var ast models.AssociateStation

	err := dao.DB.Get(&ast,
		`SELECT ast_id, ast_as_id, ast_station_type, ast_start_time, ast_end_time
        FROM associate_stations
        WHERE ast_as_id = $1 AND ast_end_time IS NULL;`, as_id)
	return ast, err
}

// CreateAssociateStation creates a new AssociateStation record after marking any existing record as completed
func (dao *AssociateDAO) CreateAssociateStation(as_id int, station_type string) (string, error) {
	var token string
	tx, err := dao.DB.Beginx()
	if err != nil {
		tx.Rollback()
		return token, err
	}

	//set currently active record to inactive
	_, err = tx.Exec(`UPDATE associate_stations SET ast_end_time = now() WHERE ast_as_id = $1 AND ast_end_time IS NULL;`, as_id)
	if err != nil {
		tx.Rollback()
		return token, err
	}

	//create new record
	//TODO extract and return api token
	err = tx.Get(&token, `INSERT INTO associate_stations (ast_as_id, ast_station_type, ast_start_time)
        VALUES ($1, $2, now())
        RETURNING ast_api_token;`, as_id, station_type)
	if err != nil {
		tx.Rollback()
		return token, err
	}

	err = tx.Commit()
	return token, err
}

// EndAssociateSession marks an associate session as ended based on the associate id and auth token
func (dao *AssociateDAO) EndAssociateSession(as_id int, token string) error {
	result, err := dao.DB.Exec(`UPDATE associate_stations SET ast_end_time = now()
        WHERE ast_as_id = $1 AND ast_end_time IS NULL AND ast_api_token = $2;`, as_id, token)

	if err != nil {
		return err
	} else if rows, _ := result.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
