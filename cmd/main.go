package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
	"wbL0"
	"wbL0/nats-streaming"
	"wbL0/pkg/handler"
	"wbL0/pkg/repository"
)

func main() {
	if err := initConfig(); err != nil {
		logrus.Fatalf("init config err: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("load .env file err: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   viper.GetString("db.DBName"),
		SSLMode:  viper.GetString("db.SSLMode"),
	})
	if err != nil {
		logrus.Fatalf("init db err: %s", err.Error())
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	repos := repository.NewRepository(db)

	ordersFromDB, err := repos.GetAllData()
	if err != nil {
		logrus.Fatalf("get orders err: %s", err.Error())
	}

	cache := wbL0.NewCache()
	cache.RestoreFromDB(ordersFromDB)
	logrus.Info("cache download was successful")

	go func() {
		err = nats_streaming.ConnectingNats(db, cache)
		if err != nil {
			logrus.Fatalf("nats streaming err: %s", err.Error())
		}
	}()
	logrus.Info("message waiting...")

	router.Post("/add", handler.NewOrder(db, cache))
	router.Get("/search", handler.SearchOrder(cache))

	srv := new(wbL0.Server)
	go func() {
		err := srv.Run(viper.GetString("port"), router)
		if err != nil {
			logrus.Fatalf("run server error: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Fatalf("shutdown server error: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Fatalf("close db err: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	return viper.ReadInConfig()
}
