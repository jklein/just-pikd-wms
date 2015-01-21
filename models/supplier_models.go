// Copyright G2G Market Inc, 2015

package models

import (
	"time"
)

type Supplier struct {
	Id   int    `json:"su_id"`
	Name string `json:"su_name"`
}

type SupplierShipment struct {
	Id               int        `json:"shi_id"`
	ShipmentCode     string     `json:"shi_shipment_code"`
	SpoId            int        `json:"shi_spo_id"`
	SuId             int        `json:"shi_su_id"`
	PromisedDelivery *time.Time `json:"shi_promised_delivery"`
	ActualDelivery   *time.Time `json:"shi_actual_delivery"`
}
