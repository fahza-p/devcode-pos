package models

import "time"

type Products struct {
	ProductId         uint32     `gorm:"type:int(10) unsigned auto_increment; primaryKey" json:"productId"`
	ProductCategoryId uint32     `json:"category_id"`
	ProductDiscountId uint32     `json:"-"`
	ProductCategory   Categories `gorm:"foreignKey:ProductCategoryId" json:"category"`
	ProductDiscount   Discounts  `gorm:"foreignKey:ProductDiscountId" json:"discount"`
	ProductSku        string     `gorm:"type:varchar(100); unique; index" json:"sku"`
	ProductName       string     `gorm:"type:varchar(255)" json:"name"`
	ProductStock      int64      `gorm:"type:int(10)" json:"stock"`
	ProductPrice      float32    `gorm:"type:decimal(10,0) unsigned" json:"price"`
	ProductImage      string     `gorm:"type:varchar(500)" json:"image"`
	CreatedAt         time.Time  `gorm:"type:datetime; autoCreateTime" json:"createdAt"`
	UpdatedAt         time.Time  `gorm:"type:datetime; autoUpdateTime" json:"updatedAt"`
}
