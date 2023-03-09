package ledgerwriter

import (
	"github.com/ServiceWeaver/weaver/examples/bankofanthos/model"
	"github.com/jinzhu/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

func newTransactionRepository(databaseURI string) (*transactionRepository, error) {
	db, err := gorm.Open("postgres", databaseURI)
	if err != nil {
		return nil, err
	}
	return &transactionRepository{
		db: db,
	}, nil
}

func (r *transactionRepository) Save(transaction *model.Transaction) error {
	return r.db.Create(transaction).Error
}
