package price

import (
	"coinbase/internal/models"
	"context"
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
	uc.cache.WriteTick(tick.ProductID, tick)
}

func (uc *UseCase) FlushTicks(ctx context.Context) error {
	ticks := uc.cache.ReadAndDeleteAllTicks()
	return uc.tickRepo.CreateTicks(ctx, ticks)
}
