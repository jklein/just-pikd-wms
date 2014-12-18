package models

import (
	"gopkg.in/guregu/null.v2"
	"time"
)

type StockingPurchaseOrder struct {
	Id            int                            `db:"stocking_purchase_order_id" json:"stocking_purchase_order_id"`
	Status        string                         `db:"status" json:"status"`
	SupplierId    int                            `db:"supplier_id" json:"supplier_id"`
	DateOrdered   *time.Time                     `db:"date_ordered" json:"date_ordered"`
	DateConfirmed *time.Time                     `db:"date_confirmed" json:"date_confirmed"`
	DateShipped   *time.Time                     `db:"date_shipped" json:"date_shipped"`
	DateArrived   *time.Time                     `db:"date_arrived" json:"date_arrived"`
	Products      []StockingPurchaseOrderProduct `json:"products"`
}

//note we don't map stocking_purchase_order_id in this struct because it is used as an element of the StockingPurchaseOrder in go anyway
//BUT how do we deal with that for creating new ones with bindvars? maybe pass in a separate thing? ugh
type StockingPurchaseOrderProduct struct {
	Id                      int         `db:"stocking_purchase_order_product_id" json:"stocking_purchase_order_product_id"`
	StockingPurchaseOrderId int         `db:"stocking_purchase_order_id"` //note - not exported to JSON since it will be 0 when loaded from DB
	Sku                     string      `db:"sku" json:"sku"`
	RequestedQty            int         `db:"requested_qty" json:"requested_qty"`
	ConfirmedQty            null.Int    `db:"confirmed_qty" json:"confirmed_qty"`
	WholesaleCost           null.Float  `db:"wholesale_cost" json:"wholesale_cost"`
	ExpirationClass         null.String `db:"expiration_class" json:"expiration_class"`
}
