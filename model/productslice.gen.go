package model

const TableNameProductslice = "productslice"

// Productslice mapped from table <productslice>
type Productslice struct {
	ID        string  `gorm:"column:id;primaryKey;default:gen_random_uuid()" json:"id"`
	ProductID string  `gorm:"column:product_id;not null" json:"product_id"`
	Slice     string  `gorm:"column:slice;not null" json:"slice"`
	Price     float64 `gorm:"column:price;not null" json:"price"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product"`
}

// TableName Productslice's table name
func (*Productslice) TableName() string {
	return TableNameProductslice
}
