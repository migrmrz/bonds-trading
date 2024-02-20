package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"trading.com/mx/orders/internal/config"
	"trading.com/mx/orders/internal/redis"
	"trading.com/mx/orders/internal/store"
	"trading.com/mx/orders/pkg/handlers/rest"
	"trading.com/mx/orders/pkg/service"
)

func main() {
	configPath := flag.String("conf", "./internal/config", "directory where config file is located")
	flag.Parse()

	// get config
	config, err := config.GetConfig(*configPath)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("error reading config")
	}

	db, err := initDB(config.DBStore.GetConnectionInfo())
	if err != nil {
		logrus.Fatalln("error initializing database:", err.Error())
	}

	defer db.Close()

	usersStore, storeErr := store.NewStorage(db, config.Queries, &config.DBStore)
	if storeErr != nil {
		logrus.Fatalf("unable to init database store: %s", storeErr.Error())
	}

	redisStore, redisErr := redis.NewStorage(config.Redis)
	if redisErr != nil {
		logrus.Fatalf("unable to init redis store: %s", redisErr.Error())
	}

	// init service.
	srv := service.New(usersStore, redisStore)

	handler := muxHandlers(srv)

	// run http service & wait.
	runService(handler)

}

func muxHandlers(service rest.OrdersHandler) *http.ServeMux {
	handler := rest.MakeHTTPHandlers(service)

	// healthcheck for server
	healthcheckHandler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	)

	// init handler with api and healthcheck
	muxHandler := http.NewServeMux()
	muxHandler.Handle("/healthcheck", healthcheckHandler)
	muxHandler.Handle("/", handler)

	return muxHandler
}

func runService(handler *http.ServeMux) {
	// init server
	server := &http.Server{
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         fmt.Sprintf(":%d", 8001),
	}

	errc := make(chan error, 2)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errc <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logrus.WithField("port", 8001).Info("listening")
		errc <- server.ListenAndServe()
	}()

	logrus.WithFields(logrus.Fields{
		"port": 8001,
	}).Info("orders-api initialized...")

	logrus.WithFields(logrus.Fields{
		"reason": <-errc,
	}).Info("terminated")
}

func initDB(connectionInfo string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
