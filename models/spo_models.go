// Copyright G2G Market Inc, 2015

package models

import (
	"gopkg.in/guregu/null.v2"
	"time"
)

type StockingPurchaseOrder struct {
	Id            int                            `json:"spo_id"`
	Status        string                         `json:"spo_status"`
	SuId          int                            `json:"spo_su_id"`
	DateOrdered   *time.Time                     `json:"spo_date_ordered"`
	DateConfirmed *time.Time                     `json:"spo_date_confirmed"`
	DateShipped   *time.Time                     `json:"spo_date_shipped"`
	DateArrived   *time.Time                     `json:"spo_date_arrived"`
	Products      []StockingPurchaseOrderProduct `json:"products"`
}

type StockingPurchaseOrderProduct struct {
	Id               int         `json:"spop_id"`
	SpoId            int         `json:"spop_spo_id"`
	PrSku            string      `json:"spop_pr_sku"`
	RclId            null.String `json:"spop_rcl_id"`
	Status           string      `json:"spop_status"`
	RequestedQty     int         `json:"spop_requested_qty"`
	ConfirmedQty     null.Int    `json:"spop_confirmed_qty"`
	ReceivedQty      null.Int    `json:"spop_received_qty"`
	CaseUpc          null.String `json:"spop_case_upc"`
	UnitsPerCase     null.Int    `json:"spop_units_per_case"`
	RequestedCaseQty null.Int    `json:"spop_requested_case_qty"`
	ConfirmedCaseQty null.Int    `json:"spop_confirmed_case_qty"`
	ReceivedCaseQty  null.Int    `json:"spop_received_case_qty"`
	CaseLength       null.Float  `json:"spop_case_length"`
	CaseWidth        null.Float  `json:"spop_case_width"`
	CaseHeight       null.Float  `json:"spop_case_height"`
	CaseWeight       null.Float  `json:"spop_case_weight"`
	ExpectedArrival  *time.Time  `json:"spop_expected_arrival"`
	ActualArrival    *time.Time  `json:"spop_actual_arrival"`
	WholesaleCost    null.Int    `json:"spop_wholesale_cost"`
	ExpirationClass  null.String `json:"spop_expiration_class"`
	MaId             int         `json:"spop_ma_id"`
}
