// Copyright G2G Market Inc, 2015

package daos

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"just-pikd-wms/models"
)

type StockingPurchaseOrderDAO struct {
	*sqlx.DB
}

// JoinedSpo represents the join of the two tables, stocking_purchase_orders and stocking_purchase_order_products
// sqlx will automatically map the columns to the right embedded struct as long as names aren't ambiguous
// column names that are ambiguous have to be aliased both in the query and in the struct tags
type JoinedSpo struct {
	models.StockingPurchaseOrder
	models.StockingPurchaseOrderProduct
}

// GetSPO retrieves an SPO object from the database based on its id, assembling models
// for the SPO and embedding models for each of its products as well
func (dao *StockingPurchaseOrderDAO) GetSPO(spo_id int) (models.StockingPurchaseOrder, error) {
	var spo models.StockingPurchaseOrder
	var rows []JoinedSpo

	//single query to get the spo and all of its products
	err := dao.DB.Select(&rows,
		`select spo_id, spo_status, spo_su_id,
        spo_date_ordered, spo_date_confirmed, spo_date_shipped, spo_date_arrived,
        spop_id, spop_spo_id, spop_pr_sku, spop_status, spop_requested_qty,
        spop_confirmed_qty, spop_received_qty, spop_case_upc, spop_units_per_case,
        spop_requested_case_qty, spop_confirmed_case_qty, spop_received_case_qty, spop_case_length,
        spop_case_width, spop_case_height, spop_case_weight, spop_expected_arrival, spop_actual_arrival,
        spop_wholesale_cost, spop_expiration_class, spop_rcl_id, spop_ma_id
        from stocking_purchase_orders
        join stocking_purchase_order_products on spop_spo_id = spo_id
        where spo_id = $1
        order by spop_id;`, spo_id)

	if err != nil {
		return spo, err
	}

	// the sqlx.Select method does not return ErrNoRows if no rows found so we have to check for it manually so the app can return 404
	if len(rows) == 0 {
		return spo, sql.ErrNoRows
	}

	//assemble the spo to return from the results
	for i, c := range rows {
		if i == 0 {
			//only need to keep first spo object since all will be the same as we're selecting by primary key
			spo = c.StockingPurchaseOrder
		}
		//append each product object to the spo object
		spo.Products = append(spo.Products, c.StockingPurchaseOrderProduct)
	}

	return spo, nil
}

// GetSPOProduct retrieves a stocking_purchase_order_product object from the database based on its id
func (dao *StockingPurchaseOrderDAO) GetSPOProduct(spo_id int, spop_id int) (models.StockingPurchaseOrderProduct, error) {
	var product models.StockingPurchaseOrderProduct

	err := dao.DB.Get(&product,
		`select spop_id, spop_spo_id, spop_pr_sku, spop_status, spop_requested_qty,
        spop_confirmed_qty, spop_received_qty, spop_case_upc, spop_units_per_case,
        spop_requested_case_qty, spop_confirmed_case_qty, spop_received_case_qty, spop_case_length,
        spop_case_width, spop_case_height, spop_case_weight, spop_expected_arrival, spop_actual_arrival,
        spop_wholesale_cost, spop_expiration_class, spop_rcl_id, spop_ma_id
        from stocking_purchase_order_products
        where spop_spo_id = $1 AND spop_id = $2;`, spo_id, spop_id)
	return product, err
}

// GetSPOs retrieves stocking purchase orders based on passed in filters, or all SPOs
func (dao *StockingPurchaseOrderDAO) GetSPOs(supplier_id int, shipment_code string) ([]models.StockingPurchaseOrder, error) {
	var spos []models.StockingPurchaseOrder

	sql_string := `select spo_id, spo_status, spo_su_id,
        spo_date_ordered, spo_date_confirmed, spo_date_shipped, spo_date_arrived,
        spop_id, spop_spo_id, spop_pr_sku, spop_status, spop_requested_qty,
        spop_confirmed_qty, spop_received_qty, spop_case_upc, spop_units_per_case,
        spop_requested_case_qty, spop_confirmed_case_qty, spop_received_case_qty, spop_case_length,
        spop_case_width, spop_case_height, spop_case_weight, spop_expected_arrival, spop_actual_arrival,
        spop_wholesale_cost, spop_expiration_class, spop_rcl_id, spop_ma_id
        from stocking_purchase_orders
        join stocking_purchase_order_products on spop_spo_id = spo_id
        `

	// use an anonymous struct for args since it's much easier to pass to a dynamically generated query and only use the required params
	args := struct {
		SupplierId   int    `json:"supplier_id"`
		ShipmentCode string `json:"shipment_code"`
	}{supplier_id, shipment_code}

	// slice of where clause conditions based on whether params are set to their 0-value or not
	var conditions []string

	if supplier_id > 0 {
		conditions = append(conditions, "spo_su_id = :supplier_id")
	}

	if len(shipment_code) > 0 {
		sql_string += "join supplier_shipments on shi_spo_id = spo_id "
		conditions = append(conditions, "shi_shipment_code = :shipment_code")
	}

	sql_string += buildWhereFromConditions(conditions) + " ORDER BY spo_id, spop_id"

	rows, err := dao.DB.NamedQuery(sql_string, args)
	if err != nil {
		return spos, err
	}
	defer rows.Close()

	//maps spo ids to slice indexes in the results
	spo_indexes := map[int]int{}
	index := -1
	for rows.Next() {
		var j JoinedSpo
		err = rows.StructScan(&j)
		if err != nil {
			return spos, err
		}

		if _, exists := spo_indexes[j.StockingPurchaseOrder.Id]; !exists {
			spos = append(spos, j.StockingPurchaseOrder)
			index += 1
			spo_indexes[j.StockingPurchaseOrder.Id] = index
		}
		spos[index].Products = append(spos[index].Products, j.StockingPurchaseOrderProduct)
	}
	err = rows.Err()

	if err == nil && len(spos) == 0 {
		return spos, sql.ErrNoRows
	}

	return spos, err
}

// CreateSPO creates a stocking purchase order object based on the pased in model, looping over the embedded
// products and creating those as well.
func (dao *StockingPurchaseOrderDAO) CreateSPO(spo_model models.StockingPurchaseOrder) (models.StockingPurchaseOrder, error) {
	// use transaction so that all inserts either succeed or fail
	tx, err := dao.DB.Beginx()
	if err != nil {
		tx.Rollback()
		return spo_model, err
	}

	// insert the base spo
	spo_model, err = dao.insertSPO(tx, spo_model)
	if err != nil {
		tx.Rollback()
		return spo_model, err
	}

	// insert each product after setting its StockingPurchaseOrderId field to the Id generated above
	for i, product := range spo_model.Products {
		product.SpoId = spo_model.Id
		result, err := dao.insertSPOProduct(tx, product)
		if err != nil {
			tx.Rollback()
			return spo_model, err
		}
		spo_model.Products[i] = result
	}

	err = tx.Commit()
	return spo_model, err
}

// CreateSPOProduct creates a stocking_purchase_order_product based on the pased in model
func (dao *StockingPurchaseOrderDAO) CreateSPOProduct(spo_product_model models.StockingPurchaseOrderProduct) (models.StockingPurchaseOrderProduct, error) {
	spo_product_model, err := dao.insertSPOProduct(dao.DB, spo_product_model)
	return spo_product_model, err
}

// insertSPO inserts a StockingPurchaseOrder model into the database and returns the model
// with its Id field updated based on the auto generated id from the database
func (dao *StockingPurchaseOrderDAO) insertSPO(e NamedExecer, spo_model models.StockingPurchaseOrder) (models.StockingPurchaseOrder, error) {
	rows, err := e.NamedQuery(
		`INSERT INTO stocking_purchase_orders (spo_status, spo_su_id, spo_date_ordered,
        spo_date_confirmed, spo_date_shipped, spo_date_arrived)
        VALUES (:spo_status, :spo_su_id, :spo_date_ordered,
        :spo_date_confirmed, :spo_date_shipped, :spo_date_arrived)
        RETURNING spo_id`,
		spo_model)

	if err != nil {
		return spo_model, err
	}

	spo_model.Id, err = extractLastInsertId(rows)
	return spo_model, err
}

// insertSPOProduct inserts a StockingPurchaseOrderProduct model into the database and returns the model
// with its Id field updated based on the auto generated id from the database
func (dao *StockingPurchaseOrderDAO) insertSPOProduct(e NamedExecer, spo_product_model models.StockingPurchaseOrderProduct) (models.StockingPurchaseOrderProduct, error) {
	rows, err := e.NamedQuery(
		`INSERT INTO stocking_purchase_order_products (spop_spo_id, spop_pr_sku, spop_status,
			spop_requested_qty, spop_confirmed_qty, spop_received_qty, spop_case_upc,
			spop_units_per_case, spop_requested_case_qty, spop_confirmed_case_qty,
			spop_received_case_qty, spop_case_length, spop_case_width, spop_case_height, spop_case_weight,
			spop_expected_arrival, spop_actual_arrival, spop_wholesale_cost, spop_expiration_class, spop_rcl_id, spop_ma_id)
        VALUES (:spop_spo_id, :spop_pr_sku, :spop_status,
        	:spop_requested_qty, :spop_confirmed_qty, :spop_received_qty, :spop_case_upc,
			:spop_units_per_case, :spop_requested_case_qty, :spop_confirmed_case_qty,
			:spop_received_case_qty, :spop_case_length, :spop_case_width, :spop_case_height, :spop_case_weight,
			:spop_expected_arrival, :spop_actual_arrival, :spop_wholesale_cost, :spop_expiration_class, :spop_rcl_id, :spop_ma_id)
        RETURNING spop_id`,
		spo_product_model)

	if err != nil {
		return spo_product_model, err
	}

	spo_product_model.Id, err = extractLastInsertId(rows)
	return spo_product_model, err
}

// UpdateSPO updates an SPO, updating only the passed-in fields
// Both the model and the raw json dictionary are required as inputs
// The model is used for properly typed and annotated fields to bind to a NamedExec sql statement
// The dict is used to generate the sql update statement based only on the fields that are actually passed in, so that we don't overwrite
// other fields with Go's 0 value for that type
func (dao *StockingPurchaseOrderDAO) UpdateSPO(spo_model models.StockingPurchaseOrder, dict map[string]interface{}) error {
	tx, err := dao.DB.Beginx()
	if err != nil {
		tx.Rollback()
		return err
	}

	// update the base SPO object if any of its fields were passed in
	stmt := buildPatchUpdate("stocking_purchase_orders", "spo_id", dict)
	err = execCheckRows(tx, stmt, spo_model)
	if err != nil {
		tx.Rollback()
		return err
	}

	// update individual product object
	if count := len(spo_model.Products); count > 0 {
		products_dict, ok := dict["products"].([]interface{})
		if !ok || len(products_dict) != count {
			tx.Rollback()
			return newInputErr("Mismatch decoding input - dict['products'] is not a slice")
		}

		// update individual products
		for i, product := range spo_model.Products {
			//redefine dict here for just that subset
			product_dict, ok := products_dict[i].(map[string]interface{})
			if !ok {
				tx.Rollback()
				return newInputErr("Mismatch decoding embedded document")
			}
			if product.SpoId > 0 && product.SpoId != spo_model.Id {
				tx.Rollback()
				return newInputErr(fmt.Sprintf("spop_spo_id does not match spo_id for product spop_id=%v", product.Id))
			}
			stmt := buildPatchUpdate("stocking_purchase_order_products", "spop_id", product_dict)

			if len(stmt) > 0 {
				if product.SpoId == 0 {
					product.SpoId = spo_model.Id
				}
				//extra check to make sure we're updating a product that is part of the correct spo_id
				stmt += " AND spop_spo_id = :spop_spo_id"

				err = execCheckRows(tx, stmt, product)
				if err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}

	err = tx.Commit()
	return err
}
