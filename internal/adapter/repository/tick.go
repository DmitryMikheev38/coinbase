package repository

import (
	"coinbase/internal/models"
	"context"
	"gorm.io/gorm"
)

type TickRepository struct {
	db *gorm.DB
}

func NewTickRepository(db *gorm.DB) *TickRepository {
	return &TickRepository{
		db: db,
	}
}

func (r *TickRepository) CreateTicks(ctx context.Context, ticks []*models.Tick) error {
	if len(ticks) != 0 {
		return r.db.WithContext(ctx).Create(ticks).Error
	}
	return nil
}
