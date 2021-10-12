package model

import (
	"errors"
	"time"
)

//ModelMap := map[string]

type User struct {
	//gorm.Model
	ID     uint   `gorm:"primary_key" rql:"filter,sort" json:",omitempty"`
	Name   string `rql:"filter,sort" json:",omitempty"`
	Age    uint   `rql:"filter,sort" json:",omitempty"`
	Agents []Agent
}

type Agent struct {
	//gorm.Model
	ID       uint      `gorm:"primary_key" rql:"filter,sort" json:",omitempty"`
	Name     string    `rql:"filter,sort" json:",omitempty"`
	CreateAt time.Time `rql:"filter,sort" json:",omitempty"`
	UserID   uint      `rql:"filter,sort" json:",omitempty"`
}

func GetObjectByName(name string) (interface{}, error) {
	switch name {
	case "User":
		return User{}, nil
	case "Agent":
		return Agent{}, nil
	default:
		return nil, errors.New("Type Not Found")
	}
}

type Params struct {
	// Limit represents the number of rows returned by the SELECT statement.
	Limit int
	// Offset specifies the offset of the first row to return. Useful for pagination.
	Offset int
	// Select contains the expression for the `SELECT` clause defined in the Query.
	Select []string
	// Sort used as a parameter for the `ORDER BY` clause. For example, "age desc, name".
	Sort string
	// FilterExp and FilterArgs come together and used as a parameters for the `WHERE` clause.
	//
	// examples:
	// 	1. Exp: "name = ?"
	//	   Args: "a8m"
	//
	//	2. Exp: "name = ? AND age >= ?"
	// 	   Args: "a8m", 22
	FilterExp  string
	FilterArgs []interface{}
}
