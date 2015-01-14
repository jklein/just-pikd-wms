// Copyright G2G Market Inc, 2015

package models

import (
	"time"
)

type Supplier struct {
	Id   int    `db:"supplier_id" json:"supplier_id"`
	Name string `db:"supplier_name" json:"supplier_name"`
}

type SupplierShipment struct {
	Id                      int        `db:"supplier_shipment_id" json:"supplier_shipment_id"`
	ShipmentId              string     `db:"shipment_id" json:"shipment_id"`
	StockingPurchaseOrderId int        `db:"stocking_purchase_order_id" json:"stocking_purchase_order_id"`
	SupplierId              int        `db:"supplier_id" json:"supplier_id"`
	PromisedDelivery        *time.Time `db:"promised_delivery" json:"promised_delivery"`
	ActualDelivery          *time.Time `db:"actual_delivery" json:"actual_delivery"`
}
