package transactionhistory

import (
	"github.com/ServiceWeaver/weaver/examples/bankofanthos/common"
	"github.com/ServiceWeaver/weaver/examples/bankofanthos/model"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// TransactionRepository is a repository for performing queries on the ledger database.
type TransactionRepository struct {
	common.LedgerReaderTransactionRepository
}

func newTransactionRepository(databaseURI string) (*TransactionRepository, error) {
	repo, err := common.NewLedgerReaderTransactionRepository(databaseURI)
	if err != nil {
		return nil, err
	}
	return &TransactionRepository{
		LedgerReaderTransactionRepository: *repo,
	}, nil
}

// findForAccount returns a list of transactions for the given account and routing numbers.
func (r *TransactionRepository) findForAccount(accountNum string, routingNum string, limit int) ([]model.Transaction, error) {
	var txns []model.Transaction
	result := r.DB.Table("transactions").Where("(from_acct=? AND from_route=?) OR (to_acct=? AND to_route=?)", accountNum, routingNum, accountNum, routingNum).Order("timestamp DESC").Find(&txns)
	if result.Error != nil {
		return nil, result.Error
	}
	return txns, nil
}
