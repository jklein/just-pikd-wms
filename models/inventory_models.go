// Copyright G2G Market Inc, 2015

package models

import (
	"fmt"
	"gopkg.in/guregu/null.v2"
	"strings"
	"time"
)

const IMAGE_URL_PREFIX string = "https://s3.amazonaws.com/g2gcdn"

type StaticInventory struct {
	Id              int         `json:"si_id"`
	StlId           string      `json:"si_stl_id"`
	PrSku           string      `json:"si_pr_sku"`
	SpopId          int         `json:"si_spop_id"`
	MaId            null.Int    `json:"si_ma_id"`
	ExpirationClass null.String `json:"si_expiration_class"`
	ExpirationDate  *time.Time  `json:"si_expiration_date"`
	TotalQty        int         `json:"si_total_qty"`
	AvailableQty    int         `json:"si_available_qty"`
	QtyOnHand       int         `json:"si_qty_on_hand"`
	ArrivalDate     *time.Time  `json:"si_arrival_date"`
	EmptiedDate     *time.Time  `json:"si_emptied_date"`
	ProductName     null.String `json:"si_product_name"`
	ProductLength   null.Float  `json:"si_product_length"`
	ProductWidth    null.Float  `json:"si_product_width"`
	ProductHeight   null.Float  `json:"si_product_height"`
	ProductWeight   null.Float  `json:"si_product_weight"`
	ThumbnailURL    string      `json:"thumbnail_url"`
}

// SetThumbnailURL computes the thumbnail URL to display to the user for
// assistance visually identifying the product.
func (s *StaticInventory) SetThumbnailURL() {
	//replace - with "" in sku, add leading 0
	sku := strings.Replace(s.PrSku, "-", "", -1)

	//result will look like https://s3.amazonaws.com/g2gcdn/68/00046000820118_200x200.jpg
	//add leading zero because they're all prefixed with that on S3 at the moment
	s.ThumbnailURL = fmt.Sprintf("%s/%d/0%s_200x200.jpg",
		IMAGE_URL_PREFIX, s.MaId.Int64, sku)
	return
}

/*
TODO: this becomes pick_task_products

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
*/

type InventoryHold struct {
	Id    int `json:"ihd_id"`
	SiId  int `json:"ihd_si_id"`
	CopId int `json:"ihd_cop_id"`
	Qty   int `json:"ihd_qty"`
}
