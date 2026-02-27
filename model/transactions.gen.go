package model

import (
	"time"

	"github.com/google/uuid"
)

const TableNameHeaderTransaction = "transactions"

// HeaderTransaction mapped from table <transactions>
type HeaderTransaction struct {
	ID                string              `gorm:"column:id;primaryKey;default:gen_random_uuid()" json:"id"`
	Customername      string              `gorm:"column:customername;not null" json:"customername"`
	Customeremail     string              `gorm:"column:customeremail;not null" json:"customeremail"`
	Address           string              `gorm:"column:address;not null" json:"address"`
	City              string              `gorm:"column:city" json:"city"`
	Postalcode        string              `gorm:"column:postalcode" json:"postalcode"`
	Totalprice        int32               `gorm:"column:totalprice;not null" json:"totalprice"`
	CreatedAt         time.Time           `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	CreatedBy         *uuid.UUID          `gorm:"column:created_by" json:"created_by"`
	DetailTransaction []DetailTransaction `gorm:"foreignKey:transaction_id" json:"detail_transaction"`
	Customer          User                `gorm:"foreignKey:created_by" json:"customer"`
}

// TableName HeaderTransaction's table name
func (*HeaderTransaction) TableName() string {
	return TableNameHeaderTransaction
}
