//TODO: using git
package main

import (
    "crypto/rand"
    "fmt"
    math_rand "math/rand"

    "gorm.io/gorm"
    "gorm.io/driver/sqlite"

    "github.com/chenfei531/query-api/model"
)

func randomStr() string {
    n := 5
    b := make([]byte, n)
    if _, err := rand.Read(b); err != nil {
        panic(err)
    }
    s := fmt.Sprintf("%X", b)
    return s
}

func randomInt() uint {
    return uint(math_rand.Intn(1000) + 1)
}

func main() {
  db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
  if err != nil {
    panic("failed to connect database")
  }

  // Migrate the schema
  db.AutoMigrate(&model.User{}, &model.Agent{})

  // Create
  for i := 0; i < 100; i++ {
      user := model.User{Name: "user_" + randomStr()}
      db.Create(&user)
      for j := 0; j < 10; j++ {
  		db.Create(&model.Agent{Name: "agent_" + randomStr(), Price: randomInt(), UserID: user.ID})
  	}
  }
}
