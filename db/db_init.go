package main

/*
import (
  "gorm.io/gorm"
  "gorm.io/driver/sqlite"

  "crypto/rand"
  "fmt"
  math_rand "math/rand"
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
  db.AutoMigrate(&User{}, &Agent{})

  // Create
  for i := 0; i < 100; i++ {
      user := User{Name: "user_" + randomStr()}
      db.Create(&user)
      for j := 0; j < 10; j++ {
  		db.Create(&Agent{Name: "agent_" + randomStr(), Price: randomInt(), UserID: user.ID})
  	}
  }
}
*/