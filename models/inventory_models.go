package models

import (
	"fmt"
	"gopkg.in/guregu/null.v2"
	"strings"
	"time"
)

const IMAGE_URL_PREFIX string = "https://s3.amazonaws.com/g2gcdn"

type StaticInventory struct {
	Id                             int         `db:"static_inventory_id" json:"static_inventory_id"`
	StockingLocationId             string      `db:"stocking_location_id" json:"stocking_location_id"`
	Sku                            string      `db:"sku" json:"sku"`
	StockingPurchaseOrderProductId int         `db:"stocking_purchase_order_product_id" json:"stocking_purchase_order_product_id"`
	ExpirationClass                null.String `db:"expiration_class" json:"expiration_class"`
	ExpirationDate                 *time.Time  `db:"expiration_date" json:"expiration_date"`
	TotalQty                       int         `db:"total_qty" json:"total_qty"`
	AvailableQty                   int         `db:"available_qty" json:"available_qty"`
	ArrivalDate                    *time.Time  `db:"arrival_date" json:"arrival_date"`
	EmptiedDate                    *time.Time  `db:"emptied_date" json:"emptied_date"`
	ManufacturerId                 null.Int    `db:"manufacturer_id" json:"manufacturer_id"`
	Name                           null.String `db:"name" json:"name"`
	Length                         null.Float  `db:"length" json:"length"`
	Width                          null.Float  `db:"width" json:"width"`
	Height                         null.Float  `db:"height" json:"height"`
	Weight                         null.Float  `db:"weight" json:"weight"`
	ThumbnailURL                   string      `json:"thumbnail_url"`
}

//thumbanil URL to display in pick app
func (s *StaticInventory) SetThumbnailURL() {
	//https://s3.amazonaws.com/g2gcdn/68/00046000820118_200x200.jpg
	//replace - with "" in sku, add leading 0
	sku := strings.Replace(s.Sku, "-", "", -1)

	//result will look like https://s3.amazonaws.com/g2gcdn/68/00046000820118_200x200.jpg
	//add leading zero because they're all prefixed with that on S3 at the moment
	s.ThumbnailURL = fmt.Sprintf("%s/%d/0%s_200x200.jpg",
		IMAGE_URL_PREFIX, s.ManufacturerId.Int64, sku)
	return
}

type OutboundInventory struct {
	Id                     int    `db:"outbound_inventory_id" json:"outbound_inventory_id"`
	CustomerOrderProductId int    `db:"customer_order_product_id" json:"customer_order_product_id"`
	Sku                    string `db:"sku" json:"sku"`
	PickContainerId        string `db:"pick_container_id" json:"pick_container_id"`
	StockingLocationId     string `db:"stocking_location_id" json:"stocking_location_id"`
	StaticInventoryId      int    `db:"static_inventory_id" json:"static_inventory_id"`
	Qty                    int    `db:"qty" json:"qty"`
	OutboundInventoryType  string `db:"outbound_inventory_type" json:"outbound_inventory_type"`
	Status                 string `db:"status" json:"status"`
}

type InventoryHold struct {
	Id                     int `db:"inventory_hold_id" json:"inventory_hold_id"`
	StaticInventoryId      int `db:"static_inventory_id" json:"static_inventory_id"`
	CustomerOrderProductId int `db:"customer_order_product_id" json:"customer_order_product_id"`
	Qty                    int `db:"qty" json:"qty"`
}
