package entities

type ProductMedia struct {
	Base

	Link      string `gorm:"column:LINK"`
	ProductID uint   `gorm:"column:PRODUCTID"`
}

func (ProductMedia) TableName() string {
	return "PRODUCTMEDIA"
}
