package models

type Supplier struct {
	Id   int    `db:"supplier_id" json:"supplier_id"`
	Name string `db:"supplier_name" json:"supplier_name"`
}

type SupplierShipment struct {
	Id                      int    `db:"supplier_shipment_id" json:"supplier_shipment_id"`
	ShipmentId              string `db:"shipment_id" json:"shipment_id"`
	StockingPurchaseOrderId int    `db:"stocking_purchase_order_id" json:"stocking_purchase_order_id"`
	SupplierId              int    `db:"supplier_id" json:"supplier_id"`
}
