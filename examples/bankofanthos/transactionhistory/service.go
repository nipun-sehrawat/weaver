package transactionhistory

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
	// GetTransactions returns all the transactions of an account.
	GetTransactions(ctx context.Context, accountID string) ([]model.Transaction, error)
}

type config struct {
	LocalRoutingNum string `toml:"local_routing_num"`
	DataSourceURL   string `toml:"data_source_url"`
	HistoryLimit    int    `toml:"history_limit"`
	CacheSize       int    `toml:"cache_size"`
	CacheMinutes    int    `toml:"cache_minutes"`
}

type impl struct {
	weaver.Implements[T]
	weaver.WithConfig[config]

	txnRepo      *TransactionRepository
	txnCache     *TransactionCache
	ledgerReader *common.LedgerReader
}

func (i *impl) ProcessTransaction(transaction model.Transaction) {
	fromID := transaction.FromAccountNum
	fromRoutingNum := transaction.FromRoutingNum
	toID := transaction.ToAccountNum
	toRouting := transaction.ToRoutingNum
	if fromRoutingNum == i.Config().LocalRoutingNum {
		if _, ok := i.txnCache.c.GetIfPresent(fromID); ok {
			i.processTransactionForAcct(fromID, transaction)
		}
	}
	if toRouting == i.Config().LocalRoutingNum {
		if _, ok := i.txnCache.c.GetIfPresent(toID); ok {
			i.processTransactionForAcct(toID, transaction)
		}
	}
}

func (i *impl) processTransactionForAcct(accountID string, transaction model.Transaction) {
	i.Logger().Debug("Processing transaction", "accountID", accountID, "transaction", transaction)
	got, err := i.txnCache.c.Get(accountID)
	if err != nil {
		i.Logger().Error("processTransactionForAcct failed", err)
		return
	}
	txns := got.([]model.Transaction)
	txns = append([]model.Transaction{transaction}, txns...)
	// Drop old transactions.
	if len(txns) > i.Config().HistoryLimit {
		i.Logger().Debug("Hit transaction caching limit, dropping old transactions", "dropped", len(txns)-i.Config().HistoryLimit)
		txns = txns[:i.Config().HistoryLimit]
	}
	i.txnCache.c.Put(accountID, txns)
}

func (i *impl) Init(context.Context) error {
	var err error
	i.txnRepo, err = newTransactionRepository(i.Config().DataSourceURL)
	if err != nil {
		return err
	}
	i.txnCache = newTransactionCache(i.txnRepo, i.Config().CacheSize, i.Config().CacheMinutes, i.Config().LocalRoutingNum, i.Config().HistoryLimit)
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

func (i *impl) GetTransactions(ctx context.Context, accountID string) ([]model.Transaction, error) {
	// Load from cache.
	got, err := i.txnCache.c.Get(accountID)
	if err != nil {
		return nil, err
	}
	return got.([]model.Transaction), nil
}
