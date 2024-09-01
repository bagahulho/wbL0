package tables

type Order struct {
	OrderUID    string `json:"order_uid,omitempty"`
	TrackNumber string `json:"track_number,omitempty"`
	Entry       string `json:"entry,omitempty"`

	Delivery Delivery `json:"delivery,omitempty"`
	Payment  Payment  `json:"payment,omitempty"`
	Items    []Item   `json:"items,omitempty"`

	Locale            string `json:"locale,omitempty"`
	InternalSignature string `json:"internal_signature,omitempty"`
	CustomerID        string `json:"customer_id,omitempty"`
	DeliveryService   string `json:"delivery_service,omitempty"`
	Shardkey          string `json:"shardkey,omitempty"`
	SmID              int    `json:"sm_id,omitempty"`
	DateCreated       string `json:"date_created,omitempty"`
	OofShard          string `json:"oof_shard,omitempty"`
}
