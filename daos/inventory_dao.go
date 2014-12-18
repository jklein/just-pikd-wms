package daos

import (
	"github.com/jmoiron/sqlx"
	"just-pikd-wms/models"
)

type Inventory_DAO struct {
	*sqlx.DB
}

func (dao *Inventory_DAO) Get_Inbound(inbound_inventory_id int) (models.InboundInventory, error) {
	var inbound models.InboundInventory

	err := dao.DB.Get(&inbound,
		`SELECT inbound_inventory_id, stocking_purchase_order_product_id, sku, expected_arrival, actual_arrival,
        confirmed_qty, received_qty, status, expiration_class, stocking_location_id
        FROM inbound_inventory
        WHERE inbound_inventory_id = $1;`, inbound_inventory_id)
	return inbound, err
}

//create an inbound inventory record and return the model created, including the newly populated ID field
func (dao *Inventory_DAO) Create_Inbound(inbound_model models.InboundInventory) (models.InboundInventory, error) {
	//insert using NamedQuery instead of NamedExec due to the need of getting the last inserted ID back
	//see https://github.com/lib/pq/issues/24
	rows, err := dao.DB.NamedQuery(
		`INSERT INTO inbound_inventory (stocking_purchase_order_product_id, sku,
        expected_arrival, actual_arrival, confirmed_qty, status, expiration_class)
        VALUES (:stocking_purchase_order_product_id, :sku,
        :expected_arrival, :actual_arrival, :confirmed_qty, :status, :expiration_class)
        RETURNING inbound_inventory_id`,
		inbound_model)
	if err != nil {
		return inbound_model, err
	}
	defer rows.Close()
	//get the inserted ID from the rowset, which will only ever be one row
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return inbound_model, err
		}
		inbound_model.Id = id
	}
	err = rows.Err()
	return inbound_model, err
}

func (dao *Inventory_DAO) Get_Static(static_inventory_id int) (models.StaticInventory, error) {
	var static models.StaticInventory

	err := dao.DB.Get(&static,
		`SELECT static_inventory_id, stocking_location_id, sku, inbound_inventory_id, expiration_class,
        expiration_date, total_qty, available_qty, arrival_date, wholesale_cost,
        manufacturer_id, name
        FROM static_inventory
        WHERE static_inventory_id = $1;`, static_inventory_id)
	return static, err
}

func (dao *Inventory_DAO) Create_Static(static_model models.StaticInventory) (models.StaticInventory, error) {
	//insert using NamedQuery instead of NamedExec due to the need of getting the last inserted ID back
	//see https://github.com/lib/pq/issues/24
	rows, err := dao.DB.NamedQuery(
		`INSERT INTO static_inventory (stocking_location_id, sku, inbound_inventory_id, expiration_class,
        expiration_date, total_qty, available_qty, arrival_date, wholesale_cost,
        manufacturer_id, name)
        VALUES (:stocking_location_id, :sku, :inbound_inventory_id, :expiration_class,
        :expiration_date, :total_qty, :available_qty, :arrival_date, :wholesale_cost,
        :manufacturer_id, :name)
        RETURNING static_inventory_id`,
		static_model)
	if err != nil {
		return static_model, err
	}
	defer rows.Close()
	//get the inserted ID from the rowset, which will only ever be one row
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return static_model, err
		}
		static_model.Id = id
	}
	err = rows.Err()
	return static_model, err
}
