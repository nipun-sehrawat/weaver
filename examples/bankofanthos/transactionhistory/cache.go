package transactionhistory

import (
	"time"

	cache "github.com/goburrow/cache"
)

// TransactionCache maintains a cache of transaction retrieved from the ledger database.
type TransactionCache struct {
	c cache.LoadingCache
}

func newTransactionCache(txnRepo *TransactionRepository,
	cacheSize int,
	expireMinutes int,
	localRoutingNum string, historyLimit int) *TransactionCache {
	load := func(accountID cache.Key) (cache.Value, error) {
		return txnRepo.findForAccount(accountID.(string), localRoutingNum, historyLimit)
	}

	return &TransactionCache{
		c: cache.NewLoadingCache(
			load,
			cache.WithMaximumSize(cacheSize),
			cache.WithExpireAfterWrite(time.Duration(expireMinutes)*time.Minute),
		),
	}
}
