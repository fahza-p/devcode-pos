package models

import "time"

type Orders struct {
	OrderId          uint32    `gorm:"type:int(10) unsigned auto_increment; primaryKey" json:"orderId"`
	OrderCashiersId  uint32    `gorm:"type:int(10) unsigned; index" json:"cashiersId"`
	Cashiers         Cashiers  `gorm:"foreignKey:OrderCashiersId" json:"cashier"`
	OrderPaymentId   uint32    `gorm:"type:int(10) unsigned; index" json:"paymentTypesId"`
	PaymentType      Payments  `gorm:"foreignKey:OrderPaymentId" json:"payment_type"`
	OrderRecipeId    string    `gorm:"type:varchar(10); index" json:"recipeId"`
	OrderTotalPrice  float32   `gorm:"type:decimal(10,0)" json:"totalPrice"`
	OrderTotalPaid   float32   `gorm:"type:decimal(10,0)" json:"totalPaid"`
	OrderTotalReturn float32   `gorm:"type:decimal(10,0)" json:"totalReturn"`
	OrderIsDownload  bool      `json:"isDownload"`
	CreatedAt        time.Time `gorm:"type:datetime; autoCreateTime" json:"createdAt"`
	UpdatedAt        time.Time `gorm:"type:datetime; autoUpdateTime" json:"updatedAt"`
}

type OrderDetails struct {
	DetailId                  uint32  `gorm:"type:int(10) unsigned auto_increment; primaryKey" json:"-"`
	DetailOrderId             uint32  `gorm:"type:int(10) unsigned; index" json:"orderId"`
	DetailProductId           uint32  `gorm:"type:int(10) unsigned; index" json:"productId"`
	DetailProductName         string  `gorm:"type:varchar(255); index" json:"name"`
	DetailDiscountId          uint32  `gorm:"type:int(10) unsigned" json:"discountsId"`
	DetailDiscountQty         int32   `gorm:"type:decimal(10,0)" json:"qty"`
	DetailDiscountPrice       float32 `gorm:"type:decimal(10,0)" json:"price"`
	DetailDiscountFinalPrice  float32 `gorm:"type:decimal(10,0)" json:"totalFinalPrice"`
	DetailDiscountNormalPrice float32 `gorm:"type:decimal(10,0)" json:"totalNormalPrice"`
}

type OrderProductDetail struct {
	DetailId        uint32   `gorm:"type:int(10) unsigned auto_increment; primaryKey" json:"-"`
	DetailOrderId   uint32   `gorm:"type:int(10) unsigned; index" json:"orderId"`
	DetailProductId uint32   `gorm:"type:int(10) unsigned; index" json:"productId"`
	DetailProducts  Products `gorm:"foreignKey:DetailProductId" json:"products"`
}

type Subtotal struct {
	ProductId         uint32    `gorm:"type:int(10) unsigned auto_increment; primaryKey" json:"productId"`
	ProductDiscountId uint32    `json:"-"`
	ProductDiscount   Discounts `gorm:"foreignKey:ProductDiscountId" json:"discount"`
	ProductName       string    `gorm:"type:varchar(255)" json:"name"`
	ProductStock      int64     `gorm:"type:int(10)" json:"stock"`
	ProductPrice      float32   `gorm:"type:decimal(10,0) unsigned" json:"price"`
	ProductImage      string    `gorm:"type:varchar(500)" json:"image"`
	Qty               int32     `json:"qty"`
	TotalNormalPrice  float32   `json:"totalNormalPrice"`
	TotalFinalPrice   float32   `json:"totalFinalPrice"`
}
