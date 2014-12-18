package models

import (
	"gopkg.in/guregu/null.v2"
	"time"
)

type InboundInventory struct {
	Id                             int         `db:"inbound_inventory_id" json:"inbound_inventory_id"`
	StockingPurchaseOrderProductId null.Int    `db:"stocking_purchase_order_product_id" json:"stocking_purchase_order_product_id"`
	Sku                            string      `db:"sku" json:"sku"`
	ExpectedArrival                *time.Time  `db:"expected_arrival" json:"expected_arrival"`
	ActualArrival                  *time.Time  `db:"actual_arrival" json:"actual_arrival"`
	ConfirmedQty                   int         `db:"confirmed_qty" json:"confirmed_qty"`
	ReceivedQty                    null.Int    `db:"received_qty" json:"received_qty"`
	Status                         string      `db:"status" json:"status"`
	ExpirationClass                null.String `db:"expiration_class" json:"expiration_class"`
	StockingLocationId             null.String `db:"stocking_location_id" json:"stocking_location_id"`
}

type StaticInventory struct {
	Id                 int         `db:"static_inventory_id" json:"static_inventory_id"`
	StockingLocationId string      `db:"stocking_location_id" json:"stocking_location_id"`
	Sku                string      `db:"sku" json:"sku"`
	InboundInventoryId int         `db:"inbound_inventory_id" json:"inbound_inventory_id"`
	ExpirationClass    null.String `db:"expiration_class" json:"expiration_class"`
	ExpirationDate     *time.Time  `db:"expiration_date" json:"expiration_date"`
	TotalQty           int         `db:"total_qty" json:"total_qty"`
	AvailableQty       null.Int    `db:"available_qty" json:"available_qty"`
	ArrivalDate        *time.Time  `db:"arrival_date" json:"arrival_date"`
	WholesaleCost      null.Float  `db:"wholesale_cost" json:"wholesale_cost"`
	ManufacturerId     null.Int    `db:"manufacturer_id" json:"manufacturer_id"`
	Name               null.String `db:"name" json:"name"`
}
