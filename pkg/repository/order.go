package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"wbL0"
	"wbL0/tables"
)

func InsertOrder(db *sqlx.DB, order tables.Order, c *wbL0.Cache) error {
	tx, err := db.Beginx()
	if err != nil {
		logrus.Error("internal fail", "error", err)

		return err
	}
	defer tx.Rollback()

	orderUid := order.OrderUID

	_, err = tx.Exec(`INSERT INTO orders(order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
                      VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
		order.CustomerID, order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OofShard)

	if err != nil {
		logrus.Error("failed to insert data", "error", err)

		return err
	}

	err = order.Delivery.Insert(tx, orderUid)
	if err != nil {
		logrus.Error("failed to insert data", "error", err)

		return err
	}

	err = order.Payment.Insert(tx, orderUid)
	if err != nil {
		logrus.Error("failed to insert data", "error", err)

		return err
	}

	for _, item := range order.Items {
		err = item.Insert(tx, orderUid)
		if err != nil {
			logrus.Error("failed to insert data", "error", err)

			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		logrus.Error("failed to commit transaction", "error", err)

		return err
	}
	c.SetOrder(order)

	return nil
}

func (r *Repository) GetItems(orderUid string) ([]tables.Item, error) {
	req := `SELECT
    chrt_id,
    track_number,
    price,
    rid,
    name,
    sale,
    size,
    total_price,
    nm_id,
    brand,
    status
    FROM items 
    WHERE order_uid = $1`
	rows, err := r.db.Query(req, orderUid)
	if err != nil {
		return nil, fmt.Errorf("error getting items: %w", err)
	}
	defer rows.Close()

	var items []tables.Item
	for rows.Next() {
		var item tables.Item

		err = rows.Scan(
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting items: %w", err)
		}

		items = append(items, item)
	}

	return items, nil
}

func (r *Repository) GetAllData() ([]tables.Order, error) {
	req := `
        SELECT
            o.order_uid,
            o.track_number,
            o.entry,
            o.locale,
            o.internal_signature,
            o.customer_id,
            o.delivery_service,
            o.shardkey,
            o.sm_id,
            o.date_created,
            o.oof_shard,
            d.name,
            d.phone,
            d.zip,
            d.city,
            d.address,
            d.region,
            d.email,
            p.transaction,
            p.request_id,
            p.currency,
            p.provider,
            p.amount,
            p.payment_dt,
            p.bank,
            p.delivery_cost,
            p.goods_total,
            p.custom_fee
        FROM
            orders o
        LEFT JOIN
            delivery d ON o.order_uid = d.order_uid
        LEFT JOIN
            payment p ON o.order_uid = p.order_uid
    `
	rows, err := r.db.Query(req)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []tables.Order
	for rows.Next() {
		var order tables.Order

		err = rows.Scan(
			&order.OrderUID,
			&order.TrackNumber,
			&order.Entry,
			&order.Locale,
			&order.InternalSignature,
			&order.CustomerID,
			&order.DeliveryService,
			&order.Shardkey,
			&order.SmID,
			&order.DateCreated,
			&order.OofShard,
			&order.Delivery.Name,
			&order.Delivery.Phone,
			&order.Delivery.Zip,
			&order.Delivery.City,
			&order.Delivery.Address,
			&order.Delivery.Region,
			&order.Delivery.Email,
			&order.Payment.Transaction,
			&order.Payment.RequestID,
			&order.Payment.Currency,
			&order.Payment.Provider,
			&order.Payment.Amount,
			&order.Payment.PaymentDT,
			&order.Payment.Bank,
			&order.Payment.DeliveryCost,
			&order.Payment.GoodsTotal,
			&order.Payment.CustomFee,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting order: %w", err)
		}
		items, err := r.GetItems(order.OrderUID)
		if err != nil {
			return nil, err
		}
		order.Items = items

		orders = append(orders, order)
	}

	return orders, nil
}
