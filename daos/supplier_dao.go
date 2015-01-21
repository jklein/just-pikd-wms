// Copyright G2G Market Inc, 2015

package daos

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"just-pikd-wms/models"
)

// SupplierDAO is used for data access related to
// supplier objects such as suppliers and supplier_shipments
type SupplierDAO struct {
	*sqlx.DB
}

// GetReceivingLocations retrieves supplier shipments based on passed in filters, or all shipments
func (dao *SupplierDAO) GetShipments(shipment_code string, spo_id int) ([]models.SupplierShipment, error) {
	var shipments []models.SupplierShipment
	var err error

	//TODO: there should be a more efficient way to do this, possibly by using a struct for the query args and NamedQuery instead of Select
	sql_base := `SELECT shi_id, shi_shipment_code, shi_spo_id
        shi_su_id, shi_promised_delivery, shi_actual_delivery
        FROM supplier_shipments`

	//assemble dynamic where clause based on parameters. since the actual function call is different
	//depending on the params (and lib/pq complains if unused parameters are passed) that's in the switch too
	switch {
	case len(shipment_code) > 0 && spo_id > 0:
		sql_string := sql_base + " WHERE shi_shipment_code = $1 AND shi_spo_id = $2"
		err = dao.DB.Select(&shipments, sql_string, shipment_code, spo_id)
	case len(shipment_code) > 0:
		sql_string := sql_base + " WHERE shi_shipment_code = $1"
		err = dao.DB.Select(&shipments, sql_string, shipment_code)
	case spo_id > 0:
		sql_string := sql_base + " WHERE shi_spo_id = $1"
		err = dao.DB.Select(&shipments, sql_string, spo_id)
	default:
		err = dao.DB.Select(&shipments, sql_base)
	}

	if err == nil && len(shipments) == 0 {
		return shipments, sql.ErrNoRows
	}
	return shipments, err
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

//func PostShipment to insert a supplier shipment record (during purchasing)
