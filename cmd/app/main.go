package main

import (
	"coinbase/internal/adapter/cache"
	"coinbase/internal/adapter/repository"
	"coinbase/internal/adapter/ws"
	"coinbase/internal/infra"
	"coinbase/internal/usecase/price"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	cfg := infra.NewConfig()

	db, err := infra.ConnectToDB(cfg.DB)
	if err != nil {
		panic(err)
	}

	ch := cache.NewCache()
	tickRepo := repository.NewTickRepository(db)

	priceUC := price.New(ch, tickRepo)
	wsClient := ws.NewClient(cfg.CoinbaseURL, priceUC)

	ctx, cancel := context.WithCancel(context.Background())

	go flushPriceTickJobTicker(ctx, 1*time.Second, priceUC)

	go func() {
		err = wsClient.SubscribeToTicketChannels(ctx, []string{"ETH-BTC", "ETH-USD", "BTC-EUR"})
		if err != nil {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)

	<-quit
	cancel()

	time.Sleep(3 * time.Second)
	log.Println("Exit")
}

func flushPriceTickJobTicker(ctx context.Context, t time.Duration, uc *price.UseCase) {
	ticker := time.NewTicker(t)

	for {
		select {
		case <-ctx.Done():
			err := uc.FlushTicks(context.Background())
			if err != nil {
				panic(err)
			}
			return
		case t := <-ticker.C:
			err := uc.FlushTicks(ctx)
			if err != nil {
				panic(err)
			}
			fmt.Println("Tick at", t)
		}
	}
}
