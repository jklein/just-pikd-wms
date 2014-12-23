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

type StockingPurchaseOrderProduct struct {
	Id                      int         `db:"stocking_purchase_order_product_id" json:"stocking_purchase_order_product_id"`
	StockingPurchaseOrderId int         `db:"stocking_purchase_order_id" json:"-"` //note - not exported to JSON since it will be 0 when loaded from DB
	Sku                     string      `db:"sku" json:"sku"`
	Status                  string      `db:"status" json:"status"`
	RequestedQty            int         `db:"requested_qty" json:"requested_qty"`
	ConfirmedQty            null.Int    `db:"confirmed_qty" json:"confirmed_qty"`
	ReceivedQty             null.Int    `db:"received_qty" json:"received_qty"`
	CaseUpc                 null.String `db:"case_upc" json:"case_upc"`
	UnitsPerCase            null.Int    `db:"units_per_case" json:"units_per_case"`
	RequestedCaseQty        null.Int    `db:"requested_case_qty" json:"requested_case_qty"`
	ConfirmedCaseQty        null.Int    `db:"confirmed_case_qty" json:"confirmed_case_qty"`
	ReceivedCaseQty         null.Int    `db:"received_case_qty" json:"received_case_qty"`
	CaseLength              null.Float  `db:"case_length" json:"case_length"`
	CaseWidth               null.Float  `db:"case_width" json:"case_width"`
	CaseHeight              null.Float  `db:"case_height" json:"case_height"`
	CaseWeight              null.Float  `db:"case_weight" json:"case_weight"`
	ExpectedArrival         *time.Time  `db:"expected_arrival" json:"expected_arrival"`
	ActualArrival           *time.Time  `db:"actual_arrival" json:"actual_arrival"`
	WholesaleCost           null.Int    `db:"wholesale_cost" json:"wholesale_cost"`
	ExpirationClass         null.String `db:"expiration_class" json:"expiration_class"`
	ReceivingLocationId     null.String `db:"receiving_location_id" json:"receiving_location_id"`
}
