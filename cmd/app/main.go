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

	tickerJobErr := flushPriceTickJobTicker(ctx, cfg.SecInterval*time.Second, priceUC)
	tickerChanErr := wsClient.SubscribeToTicketChannels(ctx, cfg.Coins)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)

	select {
	case err = <-tickerJobErr:
		log.Println("ticker job error:", err)
	case err = <-tickerChanErr:
		log.Println("ticker channel error:", err)
	case <-quit:
	}

	cancel()

	time.Sleep(5 * time.Second)
	log.Println("Exit")
}

func flushPriceTickJobTicker(ctx context.Context, t time.Duration, uc *price.UseCase) chan error {
	ticker := time.NewTicker(t)
	errChan := make(chan error)

	go func() {
		for {
			select {
			case <-ctx.Done():
				err := uc.FlushTicks(context.Background())
				if err != nil {
					fmt.Println("flushPriceTickJobTicker: ", err)
				}
				break
			case <-ticker.C:
				err := uc.FlushTicks(ctx)
				if err != nil {
					errChan <- err
				}
				break
			}
		}
	}()

	return errChan
}
