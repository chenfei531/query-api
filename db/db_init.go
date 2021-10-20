//TODO: using git
package main

import (
	"crypto/rand"
	"fmt"
	math_rand "math/rand"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

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

func randomInt(max int) uint {
	return uint(math_rand.Intn(max) + 1)
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&model.User{}, &model.Agent{}, &model.Target{}, &model.MonitorLog{}, &model.EventLog{})

	// Create
	for i := 0; i < 10; i++ {
		user := model.User{Name: "user_" + randomStr(), Age: randomInt(50)}
		db.Create(&user)
		for j := 0; j < 2; j++ {
			createAt := time.Now()
			agent := model.Agent{Name: "agent_" + randomStr(), CreateAt: &createAt, UserID: user.ID}
			db.Create(&agent)
			for k := 0; k < 2; k++ {
				target := model.Target{Name: "target_" + randomStr(), AgentID: agent.ID}
				db.Create(&target)
				for x := 0; x < 2; x++ {
					createAt = time.Now()
					log := model.MonitorLog{Timestamp: &createAt, Cpu: randomInt(100), Mem: randomInt(100), TargetID: target.ID}
					db.Create(&log)
					event := model.EventLog{Timestamp: &createAt, Event: "event_" + randomStr(), TargetID: target.ID}
					db.Create(&event)
				}
			}
		}
	}
}
