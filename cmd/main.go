package main

import (
	"context"
	"syscall"
	"os/signal"

	log "github.com/sirupsen/logrus"

	"gitlab.com/wb-dynamics/wb-go-advert-bid-checker/internal/service"
	"gitlab.com/wb-dynamics/wb-go-advert-bid-checker/internal/config"
)

func main() {
	log.Info("starting... ")
	conf, err := config.LoadConfig(".")
	log.Info("log level: ", conf.LogLevel)
	if err != nil {
		log.Fatalf("could not read config: %s", err)
	}
	lvl, err := log.ParseLevel(conf.LogLevel)
	if err != nil {
		log.Fatalf("could not parse log level: %s", err)
	}
	log.SetLevel(lvl)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	service.Serve(ctx, conf.Port)
}
