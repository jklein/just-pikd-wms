// Copyright G2G Market Inc, 2015

package daos

import (
	"github.com/jmoiron/sqlx"
	"just-pikd-wms/models"
)

// InventoryDAO contains data access methods related to inventory data
type InventoryDAO struct {
	*sqlx.DB
}

// GetStatic retrieves a static inventory record by its primary key, si_id
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

// CreateStatic creates a new static inventory record and returns the model with the auto generated id set
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

// UpdateStatic updates a static inventory record, updating only the passed-in fields
// Both the model and the raw json dictionary are required as inputs
// The model is used for properly typed and annotated fields to bind to a NamedExec sql statement
// The dict is used to generate the sql update statement based only on the fields that are actually passed in, so that we don't overwrite
// other fields with Go's 0 value for that type
func (dao *InventoryDAO) UpdateStatic(static_model models.StaticInventory, dict map[string]interface{}) error {
	// update the base SPO object if any of its fields were passed in
	stmt := buildPatchUpdate("static_inventory", "si_id", dict)
	err := execCheckRows(dao.DB, stmt, static_model)
	return err
}
