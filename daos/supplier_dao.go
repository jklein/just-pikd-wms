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
func (dao *SupplierDAO) GetShipments(shipment_id string, stocking_purchase_order_id int) ([]models.SupplierShipment, error) {
	var shipments []models.SupplierShipment
	var err error

	sql_base := `SELECT supplier_shipment_id, shipment_id, stocking_purchase_order_id
        supplier_id, promised_delivery, actual_delivery
        FROM supplier_shipments`

	//assemble dynamic where clause based on parameters. since the actual function call is different
	//depending on the params (and lib/pq complains if unused parameters are passed) that's in the switch too
	switch {
	case len(shipment_id) > 0 && stocking_purchase_order_id > 0:
		sql_string := sql_base + " WHERE shipment_id = $1 AND stocking_purchase_order_id = $2"
		err = dao.DB.Select(&shipments, sql_string, shipment_id, stocking_purchase_order_id)
	case len(shipment_id) > 0:
		sql_string := sql_base + " WHERE shipment_id = $1"
		err = dao.DB.Select(&shipments, sql_string, shipment_id)
	case stocking_purchase_order_id > 0:
		sql_string := sql_base + " WHERE stocking_purchase_order_id = $1"
		err = dao.DB.Select(&shipments, sql_string, stocking_purchase_order_id)
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
	result, err := dao.DB.NamedExec(`UPDATE supplier_shipments
        set actual_delivery = :actual_delivery,
        last_updated = now()
        WHERE supplier_shipment_id = :supplier_shipment_id`, shipment)

	if err == nil {
		// if the update doesn't match any rows, return this so the client knows it was unsuccessful
		if rows, _ := result.RowsAffected(); rows == 0 {
			return sql.ErrNoRows
		}
	}
	return err
}

//func PostShipment to insert a supplier shipment record (during purchasing)
