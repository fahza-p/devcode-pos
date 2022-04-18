package models

import "time"

type Categories struct {
	CategoryId   uint32    `gorm:"type:int(10) unsigned auto_increment; primaryKey" json:"categoryId"`
	CategoryName string    `gorm:"type:varchar(255)" json:"name"`
	CreatedAt    time.Time `gorm:"type:datetime; autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"type:datetime; autoUpdateTime" json:"updatedAt"`
}
