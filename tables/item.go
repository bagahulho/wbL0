package tables

import "github.com/jmoiron/sqlx"

type Item struct {
	ChrtID      int    `json:"chrt_id,omitempty"`
	TrackNumber string `json:"track_number,omitempty"`
	Price       int    `json:"price,omitempty"`
	Rid         string `json:"rid,omitempty"`
	Name        string `json:"name,omitempty"`
	Sale        int    `json:"sale,omitempty"`
	Size        string `json:"size,omitempty"`
	TotalPrice  int    `json:"total_price,omitempty"`
	NmID        int    `json:"nm_id,omitempty"`
	Brand       string `json:"brand,omitempty"`
	Status      int    `json:"status,omitempty"`
}

func (i *Item) Insert(t *sqlx.Tx, orderUid string) error {
	_, err := t.Exec("INSERT INTO items(order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
		orderUid, i.ChrtID, i.TrackNumber, i.Price, i.Rid, i.Name, i.Sale, i.Size, i.TotalPrice, i.NmID, i.Brand, i.Status)
	if err != nil {
		return err
	}

	return nil
}
