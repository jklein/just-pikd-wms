// Copyright G2G Market Inc, 2015

package daos

import (
	"github.com/jmoiron/sqlx"
	"just-pikd-wms/models"
)

type InventoryDAO struct {
	*sqlx.DB
}

func (dao *InventoryDAO) GetStatic(si_id int) (models.StaticInventory, error) {
	var static models.StaticInventory

	err := dao.DB.Get(&static,
		`SELECT si_id, si_stl_id, si_pr_sku, si_spop_id, si_ma_id, si_expiration_class,
        si_expiration_date, si_total_qty, si_available_qty, si_qty_on_hand,
        si_arrival_date, si_emptied_date, si_product_name, si_product_length,
        si_product_width, si_product_height, si_product_weight
        FROM static_inventory
        WHERE si_id = $1;`, si_id)
	return static, err
}

func (dao *InventoryDAO) CreateStatic(static_model models.StaticInventory) (models.StaticInventory, error) {
	//insert using NamedQuery instead of NamedExec due to the need of getting the last inserted ID back
	//see https://github.com/lib/pq/issues/24
	rows, err := dao.DB.NamedQuery(
		`INSERT INTO static_inventory (si_stl_id, si_pr_sku, si_spop_id, si_ma_id, si_expiration_class,
        si_expiration_date, si_total_qty, si_available_qty, si_qty_on_hand,
        si_arrival_date, si_emptied_date, si_product_name, si_product_length,
        si_product_width, si_product_height, si_product_weight)
        VALUES (:si_stl_id, :si_pr_sku, :si_spop_id, :si_ma_id, :si_expiration_class,
        :si_expiration_date, :si_total_qty, :si_available_qty, :si_qty_on_hand,
        :si_arrival_date, :si_emptied_date, :si_product_name, :si_product_length,
        :si_product_width, :si_product_height, :si_product_weight)
        RETURNING si_id`,
		static_model)
	if err != nil {
		return static_model, err
	}
	id, err := extractLastInsertId(rows)
	static_model.Id = id
	return static_model, err
}
