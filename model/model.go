package model

import (
  //"gorm.io/gorm"
)

type User struct {
  //gorm.Model
  ID uint `gorm:"primary_key" rql:"filter,sort"`
  Name  string `rql:"filter,sort"`
  Agents []Agent
}

type Agent struct {
    //gorm.Model
    ID uint `gorm:"primary_key" rql:"filter,sort"`
    Name string `rql:"filter,sort"`
    Price uint `rql:"filter"`
    UserID uint `rql:"filter"`
}