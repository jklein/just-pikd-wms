// Copyright G2G Market Inc, 2015

package daos

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

// NamedExecer is a custom interface to allow our methods to accept either a sqlx.DB or a sqlx.Tx as an argument,
// both of which have these methods
type NamedExecer interface {
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
}

// buildPatchUpdate is used to build an UPDATE SQL statement
// table_name is the table to be updated
// id_col is the id_column name - used in the where clause, and ignored in the SET clause
// All other keys in the map are used to create statements for the SET clause
// Values are ignored, callers should take this string and bind a struct to it instead to protect against SQL injection
func buildPatchUpdate(table_name string, id_col string, dict map[string]interface{}) string {
	columns := []string{}
	for key, val := range dict {
		// ignore embedded JSON objects since they are handled as separate table updates
		// also ignore the id column since it's immutable
		if _, is_slice := val.([]interface{}); !is_slice && key != id_col {
			// add "column = :column" to list
			columns = append(columns, fmt.Sprintf("%s = :%s", key, key))
		}
	}

	// return empty string if no columns to update
	if len(columns) == 0 {
		return ""
	}

	return fmt.Sprintf("UPDATE %s SET %s, last_updated = now() WHERE %s = :%s",
		table_name, strings.Join(columns, ", "), id_col, id_col)
}

// extractLastInsertId is a helper function to extract the LastInsertId value from
// a *sqlx.Rows result set. It extracts the int value and helps DRY some verbose error checking
func extractLastInsertId(rows *sqlx.Rows) (int, error) {
	var id int
	defer rows.Close()
	//get the inserted ID from the rowset, which should only be one row
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return id, err
		}
	}
	err := rows.Err()
	return id, err
}

// execCheckRows is a helper function to execute a SQL statement (generally an update)
// against a db, binding a passed in struct to the statement
// It does nothing if the statement is an empty string, saving callers from having to check that,
// and returns a sql.ErrNoRows if no rows are affected (which can bubble up to a 404 for the client)
func execCheckRows(e NamedExecer, stmt string, model interface{}) error {
	if len(stmt) > 0 {
		result, err := e.NamedExec(stmt, model)
		if err != nil {
			return err
		} else if rows, _ := result.RowsAffected(); rows == 0 {
			return sql.ErrNoRows
		}
	}

	return nil
}
