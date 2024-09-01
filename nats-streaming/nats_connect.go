package nats_streaming

import (
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/nats-io/stan.go"
	"github.com/sirupsen/logrus"
	"log"
	"time"
	"wbL0"
	"wbL0/pkg/repository"
	"wbL0/tables"
)

func ConnectingNats(db *sqlx.DB, c *wbL0.Cache) error {
	fmt.Println("Установка соединения с сервером NATS Streaming")
	sc, err := stan.Connect("test-cluster", "client-test", stan.NatsURL("nats://localhost:4222"), stan.ConnectWait(time.Minute))
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println("Установка соединения с сервером NATS Streaming - успешно")
	defer sc.Close()
	var order tables.Order
	subscription, err := sc.Subscribe("orders", func(msg *stan.Msg) {
		logrus.Info("Получено новое сообщение")

		err = json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Fatalf("Ошибка парсинга JSON: %v", err)

		}

		if !isValidNewOrder(order) {
			logrus.Fatal("invalid order data")
			return
		}

		// Отправка заказа в базу данных
		err := repository.InsertOrder(db, order, c)
		if err != nil {
			log.Fatalf("Ошибка записи заказа в базу данных: %v", err)
			return
		}

	})

	if err != nil {
		log.Fatal(err)
	}
	defer subscription.Unsubscribe()
	select {}
	//go printCachePeriodically(c, time.Second*15)
}

func isValidNewOrder(order tables.Order) bool {
	if order.OrderUID == "" || order.TrackNumber == "" || order.CustomerID == "" {
		return false
	}
	if order.Delivery.Name == "" || order.Delivery.Phone == "" {
		return false
	}
	if order.Payment.Transaction == "" {
		return false
	}
	for _, item := range order.Items {
		if item.ChrtID == 0 || item.TrackNumber == "" || item.Price == 0 || item.Name == "" {
			return false
		}
	}
	return true
}
