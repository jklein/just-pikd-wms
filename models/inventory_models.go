// Copyright G2G Market Inc, 2015

package models

import (
	"gopkg.in/guregu/null.v2"
	"time"
)

type StaticInventory struct {
	Id              int         `json:"si_id"`
	StlId           string      `json:"si_stl_id"`
	PrSku           string      `json:"si_pr_sku"`
	SpopId          int         `json:"si_spop_id"`
	MaId            int         `json:"si_ma_id"`
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
	s.ThumbnailURL = ThumbnailURL(s.PrSku, s.MaId)
}

type InventoryHold struct {
	Id    int `json:"ihd_id"`
	SiId  int `json:"ihd_si_id"`
	CopId int `json:"ihd_cop_id"`
	Qty   int `json:"ihd_qty"`
}
