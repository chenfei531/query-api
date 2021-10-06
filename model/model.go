package main

import (
  "gorm.io/gorm"
  //"gorm.io/driver/sqlite"
)

type User struct {
  gorm.Model
  Name  string
  Agents []Agent
}

type Agent struct {
    gorm.Model
    Name string
    Price uint
    UserID uint
}