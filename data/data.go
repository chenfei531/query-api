package data

import (
  "gorm.io/gorm"
  "gorm.io/driver/sqlite"

  "fmt"
  "github.com/chenfei531/query-api/model"
)

func GetUserWithAgent() []model.User {
    db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    var users []model.User
    db.Limit(10).Preload("Agents").Find(&users)
    return users
}

func main() {
    /*
    db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    var users []model.User
    db.Limit(10).Preload("Agents").Find(&users)
    */
    users := GetUserWithAgent()
    //fmt.Printf("%s", users[1])
    fmt.Printf("%s \n", users[1].Agents[1].Name)
}
