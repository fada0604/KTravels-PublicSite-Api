package application

import (
	"context"
	"sync"

	"ktravels-publicsite-api/internal/features/exchange_rate/domain"
	"ktravels-publicsite-api/internal/platform/logger"
)

type RateCache struct {
	mu   sync.RWMutex
	data map[string]float64
}

func NewRateCache() *RateCache {
	return &RateCache{data: make(map[string]float64)}
}

func (c *RateCache) Set(pair string, rate float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[pair] = rate
}

func (c *RateCache) Get(pair string) (float64, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.data[pair]
	return v, ok
}

func (c *RateCache) LoadAll(ctx context.Context, repo domain.Repository, log *logger.Logger) {
	rates, err := repo.FindLatestByPair(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to preload exchange rates from MongoDB; starting with empty cache")
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	for pair, rate := range rates {
		c.data[pair] = rate
	}

	log.Info().Int("pairs", len(rates)).Msg("exchange rate cache preloaded")
}
