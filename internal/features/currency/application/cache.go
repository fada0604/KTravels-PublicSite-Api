package application

import (
	"context"
	"sync"

	"ktravels-publicsite-api/internal/features/currency/domain"
	"ktravels-publicsite-api/internal/platform/logger"
	sharedcurrency "ktravels-publicsite-api/internal/shared/currency"
)

type CurrencyCache struct {
	mu   sync.RWMutex
	data map[string]sharedcurrency.CurrencyInfo
}

func NewCurrencyCache() *CurrencyCache {
	return &CurrencyCache{data: make(map[string]sharedcurrency.CurrencyInfo)}
}

func (c *CurrencyCache) Set(code string, info sharedcurrency.CurrencyInfo) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[code] = info
}

func (c *CurrencyCache) Get(code string) (sharedcurrency.CurrencyInfo, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.data[code]
	return v, ok
}

func (c *CurrencyCache) GetAll() []sharedcurrency.CurrencyInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make([]sharedcurrency.CurrencyInfo, 0, len(c.data))
	for _, v := range c.data {
		result = append(result, v)
	}
	return result
}

func (c *CurrencyCache) LoadAll(ctx context.Context, repo domain.Repository, log *logger.Logger) {
	currencies, err := repo.FindAll(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to preload currencies from MongoDB; starting with empty cache")
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	for _, cur := range currencies {
		c.data[cur.Code] = sharedcurrency.CurrencyInfo{
			Code:        cur.Code,
			Symbol:      cur.Symbol,
			Name:        cur.Name,
			Decimals:    cur.Decimals,
			Enabled:     cur.Enabled,
			DateRestart: cur.DateRestart,
		}
	}

	log.Info().Int("currencies", len(currencies)).Msg("currency cache preloaded")
}
