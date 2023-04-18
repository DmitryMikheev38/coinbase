package price

import (
	"coinbase/internal/models"
	"context"
	"fmt"
	"time"
)

type Cache interface {
	WriteTick(key string, value *models.Tick)
	ReadAndDeleteAllTicks() []*models.Tick
}

type TickRepository interface {
	CreateTicks(ctx context.Context, models []*models.Tick) error
}

type UseCase struct {
	cache    Cache
	tickRepo TickRepository
}

func New(cache Cache, tickRepo TickRepository) *UseCase {
	return &UseCase{cache: cache, tickRepo: tickRepo}
}

func (uc *UseCase) SaveTick(tick *models.Tick) {
	if tick == nil || tick.ProductID == "" || tick.BestAsk == "" || tick.BestBid == "" {
		return
	}
	uc.cache.WriteTick(tick.ProductID, tick)
}

func (uc *UseCase) FlushTicks(ctx context.Context) error {
	ticks := uc.cache.ReadAndDeleteAllTicks()
	return uc.tickRepo.CreateTicks(ctx, ticks)
}

func (uc *UseCase) FlushPriceTickJobTicker(ctx context.Context, t time.Duration) chan error {
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
