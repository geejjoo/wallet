package main

import (
	"github.com/geejjoo/task/pkg/handler"
	"github.com/geejjoo/task/pkg/repository"
	"github.com/geejjoo/task/pkg/service"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"time"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("Failed to init configs: %s", err.Error())
	}

	cfg := &repository.Config{
		WalletTablePath:  viper.GetString("tables.wallet_table"),
		HistoryTablePath: viper.GetString("tables.history_table"),
		DriverName:       viper.GetString("db.driver_name"),
		Host:             viper.GetString("db.host"),
		Port:             viper.GetString("db.port"),
		Username:         viper.GetString("db.username"),
		Password:         viper.GetString("db.password"),
		DBName:           viper.GetString("db.dbname"),
		SSLMode:          viper.GetString("db.sslmode"),
	}

	db, err := repository.NewDB(cfg)
	if err != nil {
		logrus.Fatalf("Failed to init db: %s", err.Error())
	}

	repository := repository.NewRepository(db, cfg.WalletTablePath, cfg.HistoryTablePath)
	services := service.NewService(repository)
	handlers := handler.NewHandler(services)

	httpServer := &http.Server{
		Addr:           ":" + viper.GetString("port"),
		Handler:        handlers.InitRoutes(),
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	if err := httpServer.ListenAndServe(); err != nil {
		logrus.Fatal("Failed to while running http server %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()

}
