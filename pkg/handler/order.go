package handler

import (
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"wbL0"
	response "wbL0/api"
	"wbL0/pkg/repository"
	"wbL0/tables"
)

type TemplateData struct {
	OrderUID string
	Order    *tables.Order
}

var tmpl = template.Must(template.ParseFiles("html/search_order.html"))

func SearchOrder(c *wbL0.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderUID := r.URL.Query().Get("order_uid")
		data := TemplateData{OrderUID: orderUID}

		if orderUID != "" {
			// Получаем данные из кэша
			if order, found := c.GetOrderByUID(orderUID); found {
				data.Order = &order
			}
		}

		err := tmpl.Execute(w, data)
		if err != nil {
			logrus.Fatalf("Error executing template: %v", err)
			return
		}
	}
}

func NewOrder(s *sqlx.DB, c *wbL0.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var order tables.Order
		err := render.DecodeJSON(r.Body, &order)
		if err != nil {
			logrus.Error("failed to decode request body", "error", err)

			render.JSON(w, r, response.Error("failed to decode request"))

			return
		}

		logrus.Info("request body decoded")

		if err = validator.New().Struct(order); err != nil {
			logrus.Error("failed to validate request", "error", err)

			render.JSON(w, r, response.Error("invalid request"))

			return
		}

		err = repository.InsertOrder(s, order, c)
		if err != nil {
			logrus.Error("failed to insert order", "error", err)
		}

		render.JSON(w, r, response.OK())
	}
}
