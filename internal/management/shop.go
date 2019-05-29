package management

import "database/sql"

type Management struct {
	db *sql.DB
}

// New create new management service
func New(db *sql.DB) *Management {
	return &Management{db: db}
}

// CreateShop type
type CreateShop struct {
	Name        string
	Description string
	Photos      []string
}

func (svd *Management) CreateShop(shop *CreateShop) (shopID int64, err error) {

}
