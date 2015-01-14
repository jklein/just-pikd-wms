// Copyright G2G Market Inc, 2015

package daos

import (
	"github.com/jmoiron/sqlx"
	"just-pikd-wms/models"
)

type InventoryDAO struct {
	*sqlx.DB
}

func (dao *InventoryDAO) GetStatic(static_inventory_id int) (models.StaticInventory, error) {
	var static models.StaticInventory

	err := dao.DB.Get(&static,
		`SELECT static_inventory_id, stocking_location_id, sku, stocking_purchase_order_product_id,
        expiration_class, expiration_date, total_qty, available_qty, arrival_date, emptied_date,  manufacturer_id,
        name, length, width, height, weight
        FROM static_inventory
        WHERE static_inventory_id = $1;`, static_inventory_id)
	return static, err
}

func (dao *InventoryDAO) CreateStatic(static_model models.StaticInventory) (models.StaticInventory, error) {
	//insert using NamedQuery instead of NamedExec due to the need of getting the last inserted ID back
	//see https://github.com/lib/pq/issues/24
	rows, err := dao.DB.NamedQuery(
		`INSERT INTO static_inventory (stocking_location_id, sku, stocking_purchase_order_product_id, expiration_class,
        expiration_date, total_qty, available_qty, arrival_date, emptied_date, manufacturer_id, name, length, width, height, weight)
        VALUES (:stocking_location_id, :sku, :stocking_purchase_order_product_id, :expiration_class, :expiration_date,
            :total_qty, :available_qty, :arrival_date, :emptied_date, :manufacturer_id, :name, :length, :width, :height, :weight)
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
