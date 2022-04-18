package models

import "time"

type Cashiers struct {
	CashierId       uint32    `gorm:"type:int(10) unsigned auto_increment; primaryKey" json:"cashierId"`
	CashierName     string    `gorm:"type:varchar(255)" json:"name"`
	CashierPasscode string    `gorm:"type:varchar(100); index" json:"-"`
	CashierToken    string    `gorm:"type:text; index" json:"-"`
	CreatedAt       time.Time `gorm:"type:datetime; autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"type:datetime; autoUpdateTime" json:"updatedAt"`
}
