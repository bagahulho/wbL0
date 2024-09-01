package tables

import "github.com/jmoiron/sqlx"

type Payment struct {
	Transaction  string `json:"transaction,omitempty"`
	RequestID    string `json:"request_id,omitempty"`
	Currency     string `json:"currency,omitempty"`
	Provider     string `json:"provider,omitempty"`
	Amount       int    `json:"amount,omitempty"`
	PaymentDT    int    `json:"payment_dt,omitempty"`
	Bank         string `json:"bank,omitempty"`
	DeliveryCost int    `json:"delivery_cost,omitempty"`
	GoodsTotal   int    `json:"goods_total,omitempty"`
	CustomFee    int    `json:"custom_fee,omitempty"`
}

func (p *Payment) Insert(t *sqlx.Tx, orderUid string) error {
	_, err := t.Exec("INSERT INTO payment(order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		orderUid, p.Transaction, p.RequestID, p.Currency, p.Provider, p.Amount, p.PaymentDT, p.Bank, p.DeliveryCost, p.GoodsTotal, p.CustomFee)
	if err != nil {
		return err
	}

	return nil
}
