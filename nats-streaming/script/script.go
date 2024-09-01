package main

import (
	"encoding/json"
	"github.com/icrowley/fake"
	"github.com/nats-io/stan.go"
	"log"
	"math/rand"
	"time"
	"wbL0/tables"
)

func main() {
	sc, err := stan.Connect("test-cluster", "client", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		log.Fatal("Tut osh: ", err)
	}
	defer sc.Close()

	for i := 0; i < 5; i++ {
		order := createOrders()
		orderJSON, err := json.Marshal(order)
		if err != nil {
			log.Fatal(err)
		}
		channel := "orders"
		err = sc.Publish(channel, orderJSON)
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("Order data sent to channel")
}

func createOrders() tables.Order {
	order := tables.Order{
		OrderUID:          fake.DigitsN(13),
		TrackNumber:       fake.CharactersN(10),
		Entry:             fake.CharactersN(20),
		Locale:            "en",
		InternalSignature: fake.CharactersN(15),
		CustomerID:        fake.CharactersN(8),
		DeliveryService:   fake.Company(),
		Shardkey:          fake.CharactersN(10),
		SmID:              fake.Day(),
		DateCreated:       "2021-11-26T06:22:19Z",
		OofShard:          fake.CharactersN(5),
	}

	order.Delivery = tables.Delivery{
		Name:    fake.FullName(),
		Phone:   fake.Phone(),
		Zip:     fake.Zip(),
		City:    fake.City(),
		Address: fake.StreetAddress(),
		Region:  fake.State(),
		Email:   fake.EmailAddress(),
	}

	order.Payment = tables.Payment{
		Transaction:  fake.CharactersN(10),
		RequestID:    fake.CharactersN(8),
		Currency:     fake.CurrencyCode(),
		Provider:     fake.Company(),
		Amount:       3,
		PaymentDT:    int(time.Now().Unix()),
		Bank:         "Tinkoff",
		DeliveryCost: rand.Intn(10000),

		GoodsTotal: rand.Intn(100),
		CustomFee:  rand.Intn(10),
	}
	for i := 0; i < 3; i++ {
		item := tables.Item{
			ChrtID:      i + 1,
			TrackNumber: fake.CharactersN(8),
			Price:       rand.Intn(30000),
			Rid:         fake.CharactersN(5),
			Name:        fake.ProductName(),
			Sale:        rand.Intn(30),
			Size:        "0",
			TotalPrice:  rand.Intn(1000),
			NmID:        rand.Intn(599),
			Brand:       fake.Brand(),
			Status:      rand.Intn(400),
		}
		order.Items = append(order.Items, item)
	}
	return order
}
