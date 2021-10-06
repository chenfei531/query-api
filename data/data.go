package main

/*
import (
  "gorm.io/gorm"
  "gorm.io/driver/sqlite"

  "fmt"
)

func main() {
    db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    var users []User
    db.Limit(10).Preload("Agents").Find(&users)

    fmt.Printf("%s \n", users[1].Agents[1].Name)
}
*/