package balancereader

import (
	cache "github.com/goburrow/cache"
)

type balanceCache struct {
	c cache.LoadingCache
}

func newTransactionCache(txnRepo *TransactionRepository, expireSize int, localRoutingNum string) *balanceCache {
	load := func(accountID cache.Key) (cache.Value, error) {
		return txnRepo.findBalance(accountID.(string), localRoutingNum)
	}
	return &balanceCache{
		c: cache.NewLoadingCache(
			load,
			cache.WithMaximumSize(expireSize),
		),
	}
}
