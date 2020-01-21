package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	redis "github.com/go-redis/redis/v7"
)

func main() {

	cfg, err := initConfig()
	if err != nil {
		log.Fatalf("Initializing failed: %s", err.Error())
	}

	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{cfg.Endpoint},
	})
	if _, err := client.Ping().Result(); err != nil {
		log.Fatalf("Failed to connect to redis at %s: %s", cfg.Endpoint, err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		cancel()
	}()

	switch cfg.Mode {
	case "produce":
		switch cfg.SourceFormat {
		case "blob":
			if err := produceBlob(ctx, cfg, client); err != nil {
				log.Fatalf("produce failed: %s", err.Error())
			}
		case "csv":
			if err := produce(ctx, cfg, client); err != nil {
				log.Fatalf("produce failed: %s", err.Error())
			}
		}
	case "consume":
		if err := consume(ctx, cfg, client); err != nil {
			log.Fatalf("consume failed: %s", err.Error())
		}
	}
}
