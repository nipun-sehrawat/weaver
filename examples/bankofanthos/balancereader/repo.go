package balancereader

import (
	"github.com/ServiceWeaver/weaver/examples/bankofanthos/common"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// TransactionRepository is a repository for performing queries on the Transaction database.
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

// findBalance returns the current balance of the given account.
func (r *TransactionRepository) findBalance(accountNum, routingNum string) (int64, error) {
	sql := `
SELECT
  (
    SELECT SUM(AMOUNT) FROM TRANSACTIONS t
    WHERE (TO_ACCT = ? AND TO_ROUTE = ?)
  )
  -
  (
    SELECT COALESCE((SELECT SUM(AMOUNT) FROM TRANSACTIONS t
    WHERE (FROM_ACCT = ? AND FROM_ROUTE = ?)),0)
  )
  as balance
`
	var balance int64
	row := r.DB.Raw(sql, accountNum, routingNum, accountNum, routingNum).Row()
	row.Scan(&balance)
	return balance, nil
}
