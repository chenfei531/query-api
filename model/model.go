package model

import (
//"gorm.io/gorm"
)

type User struct {
	//gorm.Model
	ID     uint   `gorm:"primary_key" rql:"filter,sort" json:",omitempty"`
	Name   string `rql:"filter,sort" json:",omitempty"`
	Agents []Agent
}

type Agent struct {
	//gorm.Model
	ID     uint   `gorm:"primary_key" rql:"filter,sort" json:",omitempty"`
	Name   string `rql:"filter,sort" json:",omitempty"`
	Price  uint   `rql:"filter" json:",omitempty"`
	UserID uint   `rql:"filter" json:",omitempty"`
}
