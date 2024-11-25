package jobs

import (
	"context"
	"ecommerce/warehouse-service/repository"
	"ecommerce/warehouse-service/utils"
	"time"
)

type StockReleaseJob struct {
	warehouseRepo repository.WarehouseRepository
	interval      time.Duration
}

func NewStockReleaseJob(repo repository.WarehouseRepository, interval time.Duration) *StockReleaseJob {
	return &StockReleaseJob{
		warehouseRepo: repo,
		interval:      interval,
	}
}

func (j *StockReleaseJob) Start(ctx context.Context) {
	ticker := time.NewTicker(j.interval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				j.releaseExpiredStock(ctx)
			}
		}
	}()
}

func (j *StockReleaseJob) releaseExpiredStock(ctx context.Context) {
	expirationTime := time.Now().Add(-15 * time.Minute)
	utils.Logger.Infof("Running stock release job for transfers before %v", expirationTime)

	// TODO: Implement the actual stock release logic
	// 1. Find all pending transfers older than expirationTime
	// 2. Release the reserved stock back to the source warehouse
	// 3. Update transfer status to "expired"
}
