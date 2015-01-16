// Copyright G2G Market Inc, 2015

package daos

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"just-pikd-wms/models"
)

type StockingPurchaseOrderDAO struct {
	*sqlx.DB
}

func (dao *StockingPurchaseOrderDAO) GetSPO(spo_id int) (models.StockingPurchaseOrder, error) {
	var spo models.StockingPurchaseOrder

	//create an inline type to represent the join of the two tables
	//sqlx will automatically map the columns to the right embedded struct as long as names aren't ambiguous
	type JoinedSpo struct {
		models.StockingPurchaseOrder
		models.StockingPurchaseOrderProduct
	}
	rows := []JoinedSpo{}

	//single query to get the spo and all of its products
	err := dao.DB.Select(&rows,
		`select spo.stocking_purchase_order_id, spo.status, spo.supplier_id,
        spo.date_ordered, spo.date_confirmed, spo.date_shipped, spo.date_arrived,
        spop.stocking_purchase_order_product_id, spop.sku, spop.status as spopstatus, spop.requested_qty,
        spop.confirmed_qty, spop.received_qty, spop.case_upc, spop.units_per_case,
        spop.requested_case_qty, spop.confirmed_case_qty, spop.received_case_qty, spop.case_length,
        spop.case_width, spop.case_height, spop.case_weight, spop.expected_arrival, spop.actual_arrival,
        spop.wholesale_cost, spop.expiration_class, spop.receiving_location_id
        from stocking_purchase_orders spo
        join stocking_purchase_order_products spop using (stocking_purchase_order_id)
        where spo.stocking_purchase_order_id = $1;`, spo_id)

	if err != nil {
		return spo, err
	}

	// the sqlx.Select method does not return ErrNoRows if no rows found so we have to check for it manually
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

func (dao *StockingPurchaseOrderDAO) CreateSPO(spo_model models.StockingPurchaseOrder) (models.StockingPurchaseOrder, error) {
	// use transaction so that all inserts either succeed or fail
	tx, err := dao.DB.Beginx()
	if err != nil {
		tx.Rollback()
		return spo_model, err
	}

	spo_model, err = dao.insertSPO(tx, spo_model)
	if err != nil {
		tx.Rollback()
		return spo_model, err
	}

	// TODO is there a good way to do this without running multiple queries?
	// should be able to create an insert_SPO_Products method that accepts a slice of models
	// however we would probably not be able to use prepared statements in that case
	for i, product := range spo_model.Products {
		product.StockingPurchaseOrderId = spo_model.Id
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

func (dao *StockingPurchaseOrderDAO) insertSPO(tx *sqlx.Tx, spo_model models.StockingPurchaseOrder) (models.StockingPurchaseOrder, error) {
	rows, err := tx.NamedQuery(
		`INSERT INTO stocking_purchase_orders (status, supplier_id, date_ordered,
        date_confirmed, date_shipped, date_arrived)
        VALUES (:status, :supplier_id, :date_ordered,
        :date_confirmed, :date_shipped, :date_arrived)
        RETURNING stocking_purchase_order_id`,
		spo_model)
	if err != nil {
		return spo_model, err
	}
	defer rows.Close()
	//get the inserted ID from the rowset, which will only ever be one row
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return spo_model, err
		}
		spo_model.Id = id
	}
	err = rows.Err()

	return spo_model, err
}

func (dao *StockingPurchaseOrderDAO) insertSPOProduct(tx *sqlx.Tx, spo_product_model models.StockingPurchaseOrderProduct) (models.StockingPurchaseOrderProduct, error) {
	rows, err := tx.NamedQuery(
		`INSERT INTO stocking_purchase_order_products (stocking_purchase_order_id, sku, status, requested_qty, confirmed_qty, received_qty, case_upc,
			units_per_case, requested_case_qty, confirmed_case_qty, received_case_qty, case_length, case_width, case_height, case_weight,
			expected_arrival, actual_arrival, wholesale_cost, expiration_class, receiving_location_id)
        VALUES (:stocking_purchase_order_id, :sku, :spopstatus, :requested_qty, :confirmed_qty, :received_qty, :case_upc, :units_per_case,
        	:requested_case_qty, :confirmed_case_qty, :received_case_qty, :case_length, :case_width, :case_height, :case_weight,
        	:expected_arrival, :actual_arrival, :wholesale_cost, :expiration_class, :receiving_location_id)
        RETURNING stocking_purchase_order_product_id`,
		spo_product_model)
	if err != nil {
		return spo_product_model, err
	}
	defer rows.Close()
	//get the inserted ID from the rowset, which will only ever be one row
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return spo_product_model, err
		}
		spo_product_model.Id = id
	}
	err = rows.Err()

	return spo_product_model, err
}
