// Copyright G2G Market Inc, 2015

package daos

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"just-pikd-wms/models"
	"strings"
)

// SupplierDAO is used for data access related to
// supplier objects such as suppliers and supplier_shipments
type SupplierDAO struct {
	*sqlx.DB
}

// GetReceivingLocations retrieves supplier shipments based on passed in filters, or all shipments
func (dao *SupplierDAO) GetShipments(shipment_code string, spo_id int) ([]models.SupplierShipment, error) {
	var shipments []models.SupplierShipment

	// base sql string before where clause
	sql_string := `SELECT shi_id, shi_shipment_code, shi_spo_id,
        shi_su_id, shi_promised_delivery, shi_actual_delivery
        FROM supplier_shipments
        `

	// use an anonymous struct for args since it's much easier to pass to a dynamically generated query and only use the required params
	args := struct {
		ShipmentCode string `json:"shipment_code"`
		SpoId        int    `json:"spo_id"`
	}{shipment_code, spo_id}

	// slice of where clause conditions based on whether params are set to their 0-value or not
	var conditions []string

	if len(shipment_code) > 0 {
		conditions = append(conditions, "shi_shipment_code = :shipment_code")
	}

	if spo_id > 0 {
		conditions = append(conditions, "shi_spo_id = :spo_id")
	}

	sql_string += buildWhereFromConditions(conditions)

	rows, err := dao.DB.NamedQuery(sql_string, args)
	if err != nil {
		return shipments, err
	}
	defer rows.Close()

	for rows.Next() {
		var s models.SupplierShipment
		err = rows.StructScan(&s)
		if err != nil {
			return shipments, err
		}
		shipments = append(shipments, s)
	}
	err = rows.Err()

	if err == nil && len(shipments) == 0 {
		return shipments, sql.ErrNoRows
	}

	return shipments, err
}

// buildWhere builds a where clause from a slice of condition expressions, joining the conditions with AND
func buildWhereFromConditions(conditions []string) string {
	if len(conditions) > 0 {
		return "WHERE " + strings.Join(conditions, " AND ")
	}
	return ""
}

// UpdateShipment updates a supplier shipment and returns any errors it received
// The passed in model should be a shipment that already exists
// Only the mutable field - actual_delivery, is updated
// Other fields are considered immutable and are not updated.
func (dao *SupplierDAO) UpdateShipment(shipment models.SupplierShipment) error {
	stmt := `UPDATE supplier_shipments
        set shi_actual_delivery = :shi_actual_delivery,
        last_updated = now()
        WHERE shi_id = :shi_id`

	return execCheckRows(dao.DB, stmt, shipment)
}

// CreateShipment creates a new supplier shipment record, adds the auto generated id to the passed in model
// and returns it
func (dao *SupplierDAO) CreateShipment(shipment models.SupplierShipment) (models.SupplierShipment, error) {
	rows, err := dao.DB.NamedQuery(
		`INSERT INTO supplier_shipments (shi_shipment_code, shi_spo_id,
        shi_su_id, shi_promised_delivery, shi_actual_delivery)
        VALUES (:shi_shipment_code, :shi_spo_id,
        :shi_su_id, :shi_promised_delivery, :shi_actual_delivery)
        RETURNING shi_id`,
		shipment)
	if err != nil {
		return shipment, err
	}
	id, err := extractLastInsertId(rows)
	shipment.Id = id
	return shipment, err
}
