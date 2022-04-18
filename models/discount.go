package models

import "time"

type Discounts struct {
	DiscountId      uint32    `gorm:"type:int(10) unsigned auto_increment; primaryKey" json:"discountId"`
	DiscountQty     uint16    `gorm:"type:int(10) unsigned" json:"qty"`
	DiscountType    string    `gorm:"type:varchar(20); unique; index" json:"type"`
	DiscountResult  float32   `gorm:"type:decimal(10,0)" json:"result"`
	ExpiredAt       time.Time `gorm:"type:datetime" json:"expiredAt"`
	ExpiredAtFormat string    `json:"expiredAtFormat"`
	StringFormat    string    `json:"stringFormat"`
}
