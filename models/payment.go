package models

type Payments struct {
	PaymentId   uint32 `gorm:"type:int(10) unsigned auto_increment; primaryKey" json:"paymentId"`
	PaymentName string `gorm:"type:varchar(50)" json:"name"`
	PaymentType string `gorm:"type:varchar(50); index" json:"type"`
	PaymentLogo string `gorm:"type:varchar(500)" json:"logo"`
}
