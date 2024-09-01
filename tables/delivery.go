package tables

import (
	"github.com/jmoiron/sqlx"
)

type Delivery struct {
	Name    string `json:"name,omitempty"`
	Phone   string `json:"phone,omitempty"`
	Zip     string `json:"zip,omitempty"`
	City    string `json:"city,omitempty"`
	Address string `json:"address,omitempty"`
	Region  string `json:"region,omitempty"`
	Email   string `json:"email,omitempty"`
}

func (d *Delivery) Insert(t *sqlx.Tx, orderUid string) error {
	_, err := t.Exec("INSERT INTO delivery(order_uid, name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		orderUid, d.Name, d.Phone, d.Zip, d.City, d.Address, d.Region, d.Email)
	if err != nil {
		return err
	}

	return nil
}
