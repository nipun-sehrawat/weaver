package balancereader

import (
	"context"
	"errors"

	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/examples/bankofanthos/common"
	"github.com/ServiceWeaver/weaver/examples/bankofanthos/model"
)

type T interface {
	// Healthy returns the health status of this component.
	Healthy(ctx context.Context) (string, int32, error)
	// GetBalance returns the balance of an account id.
	GetBalance(ctx context.Context, accountID string) (int64, error)
}

type config struct {
	LocalRoutingNum string `toml:"local_routing_num"`
	DataSourceURL   string `toml:"data_source_url"`
}

type impl struct {
	weaver.Implements[T]
	weaver.WithConfig[config]
	txnRepo      *TransactionRepository
	balanceCache *balanceCache
	ledgerReader *common.LedgerReader
}

func (i *impl) ProcessTransaction(transaction model.Transaction) {
	fromID := transaction.FromAccountNum
	fromRoutingNum := transaction.FromRoutingNum
	toID := transaction.ToAccountNum
	toRouting := transaction.ToRoutingNum
	amount := transaction.Amount
	if fromRoutingNum == i.Config().LocalRoutingNum {
		if got, ok := i.balanceCache.c.GetIfPresent(fromID); ok {
			prevBalance := got.(int64)
			i.balanceCache.c.Put(fromID, prevBalance-int64(amount))
		}
	}
	if toRouting == i.Config().LocalRoutingNum {
		if got, ok := i.balanceCache.c.GetIfPresent(toID); ok {
			prevBalance := got.(int64)
			i.balanceCache.c.Put(toID, prevBalance+int64(amount))
		}
	}
}

func (i *impl) Init(context.Context) error {
	var err error
	i.txnRepo, err = newTransactionRepository(i.Config().DataSourceURL)
	if err != nil {
		return err
	}
	cacheSize := 1000000
	i.balanceCache = newTransactionCache(i.txnRepo, cacheSize, i.Config().LocalRoutingNum)
	i.ledgerReader = common.NewLedgerReader(i.txnRepo, i.Logger())
	i.ledgerReader.StartWithCallback(i)
	return nil
}

func (i *impl) Healthy(ctx context.Context) (string, int32, error) {
	if i.ledgerReader.IsAlive() {
		return "ok", 200, nil
	}
	err := errors.New("Ledger reader is unhealthy")
	return err.Error(), 500, err
}

func (i *impl) GetBalance(ctx context.Context, accountID string) (int64, error) {
	// Load from cache.
	got, err := i.balanceCache.c.Get(accountID)
	if err != nil {
		return 0, err
	}
	return got.(int64), nil
}
