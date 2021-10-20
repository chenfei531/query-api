package query

import (
	"errors"
	"time"
)

type User struct {
	ID     uint    `gorm:"primary_key" rql:"filter,sort" json:",omitempty"`
	Name   string  `rql:"filter,sort" json:",omitempty"`
	Age    uint    `rql:"filter,sort" json:",omitempty"`
	Agents []Agent `rqlp:"nested" json:",omitempty"`
}

type Agent struct {
	ID       uint       `gorm:"primary_key" rql:"filter,sort" json:",omitempty"`
	Name     string     `rql:"filter,sort" json:",omitempty"`
	CreateAt *time.Time `rql:"filter,sort" json:",omitempty"`
	UserID   uint       `rql:"filter,sort" json:",omitempty"`
	Targets  []Target   `rqlp:"nested" json:",omitempty"`
}

type Target struct {
	ID          uint         `gorm:"primary_key" rql:"filter,sort" json:",omitempty"`
	Name        string       `rql:"filter,sort" json:",omitempty"`
	AgentID     uint         `rql:"filter,sort" json:",omitempty"`
	MonitorLogs []MonitorLog `rqlp:"nested" json:",omitempty"`
}

type MonitorLog struct {
	ID        uint       `gorm:"primary_key" rql:"filter,sort" json:",omitempty"`
	Timestamp *time.Time `rql:"filter,sort" json:",omitempty"`
	Cpu       uint       `rql:"filter,sort" json:",omitempty"`
	Mem       uint       `rql:"filter,sort" json:",omitempty"`
	TargetID  uint       `rql:"filter,sort" json:",omitempty"`
}

func GetObjectByName(name string) (interface{}, error) {
	switch name {
	case "User":
		return User{}, nil
	case "Agent":
		return Agent{}, nil
	case "Target":
		return Target{}, nil
	case "MonitorLog":
		return MonitorLog{}, nil
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

type Node struct {
	Name       string
	PrimaryKey string
	ContainPK  bool
	Params     Params //TODO: to ptr?
	Children   []*Node
}
