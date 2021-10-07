package data

import (
    "fmt"
    "gorm.io/gorm"
    "gorm.io/driver/sqlite"

    "github.com/chenfei531/query-api/model"
)

type DataManager struct{
    db *gorm.DB
}

func NewDataManager() (*DataManager) {
    db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    return &DataManager{db}
}

func (dm DataManager) GetUserWithAgent() []model.User {
    var users []model.User
    dm.db.Limit(10).Preload("Agents").Find(&users)
    return users
}

func (dm DataManager) GetAgents(offset int, limit int) []model.Agent {
    var agents []model.Agent
    dm.db.Offset(offset).Limit(limit).Find(&agents)
    return agents
}

func (dm DataManager) GetAgentByParams(s string, fe string, fa []interface{}, offset int, limit int, sort string) []model.Agent {
    var agents []model.Agent
    error := dm.db.Select(s).
                   Where(fe, fa).
		           Offset(offset).
                   Limit(limit).
	               Order(sort).
                   Find(&agents).Error
    if error != nil {
        fmt.Printf("query error: %s \n", error)
    }
    return agents
}