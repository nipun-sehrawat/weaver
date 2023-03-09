package contacts

import (
	"github.com/ServiceWeaver/weaver"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Contact represents an account's contact details.
type Contact struct {
	weaver.AutoMarshal
	Username   string `gorm:"not null"`
	Label      string `gorm:"not null"`
	AccountNum string `gorm:"not null"`
	RoutingNum string `gorm:"not null"`
	IsExternal bool   `gorm:"not null"`
}

type contactDB struct {
	db *gorm.DB
}

func newContactDB(uri string) (*contactDB, error) {
	db, err := gorm.Open("postgres", uri)
	if err != nil {
		return nil, err
	}
	return &contactDB{
		db: db,
	}, nil
}

func (cdb *contactDB) addContact(contact Contact) error {
	return cdb.db.Create(&contact).Error
}

func (cdb *contactDB) getContacts(username string) ([]Contact, error) {
	contacts := []Contact{}
	err := cdb.db.Where("username = ?", username).Find(&contacts).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return contacts, nil
}
