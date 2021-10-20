package model

import (
	"time"
)

type User struct {
	ID     uint    `gorm:"primary_key" rql:"filter,sort" json:",omitempty"`
	Name   string  `rql:"filter,sort" json:",omitempty"`
	Age    uint    `rql:"filter,sort" json:",omitempty"`
	Agents []Agent `json:",omitempty"`
}

type Agent struct {
	ID       uint       `gorm:"primary_key" rql:"filter,sort" json:",omitempty"`
	Name     string     `rql:"filter,sort" json:",omitempty"`
	CreateAt *time.Time `rql:"filter,sort" json:",omitempty"`
	UserID   uint       `rql:"filter,sort" json:",omitempty"`
	Targets  []Target   `json:",omitempty"`
}

type Target struct {
	ID          uint         `gorm:"primary_key" rql:"filter,sort" json:",omitempty"`
	Name        string       `rql:"filter,sort" json:",omitempty"`
	AgentID     uint         `rql:"filter,sort" json:",omitempty"`
	MonitorLogs []MonitorLog `json:",omitempty"`
	EventLogs   []EventLog   `json:",omitempty"`
}

type MonitorLog struct {
	ID        uint       `gorm:"primary_key" rql:"filter,sort" json:",omitempty"`
	Timestamp *time.Time `rql:"filter,sort" json:",omitempty"`
	Cpu       uint       `rql:"filter,sort" json:",omitempty"`
	Mem       uint       `rql:"filter,sort" json:",omitempty"`
	TargetID  uint       `rql:"filter,sort" json:",omitempty"`
}

type EventLog struct {
	ID        uint       `gorm:"primary_key" rql:"filter,sort" json:",omitempty"`
	Timestamp *time.Time `rql:"filter,sort" json:",omitempty"`
	Event     string     `rql:"filter,sort" json:",omitempty"`
	TargetID  uint       `rql:"filter,sort" json:",omitempty"`
}
