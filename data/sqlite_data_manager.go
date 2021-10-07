package data

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/chenfei531/query-api/model"
)

type SqliteDataManager struct {
	db *gorm.DB
}

func NewSqliteDataManager() *SqliteDataManager {
	db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	return &SqliteDataManager{db}
}

func (dm *SqliteDataManager) GetUserById(id int) model.User {
	var user model.User
	dm.db.Preload("Agents").First(&user, id)
	return user
}

func (dm *SqliteDataManager) GetUserByParams(p *model.Params) []model.User {
	var users []model.User

	error := dm.db. //Select(p.Select).
			Where(p.FilterExp, p.FilterArgs).
			Offset(p.Offset).
			Limit(p.Limit).
			Preload("Agents").
			Order(p.Sort).
			Find(&users).Error
	if error != nil {
		fmt.Printf("query error: %s \n", error)
	}
	/*
	   // currently preload does not work with select
	   //reference: https://github.com/go-gorm/gorm/issues/4015
	   dm.db.Model(&model.User{}).Preload("Agents", func(tx *gorm.DB) *gorm.DB{
	       return tx.Select("Name").Where(p.FilterExp, p.FilterArgs).Offset(p.Offset).Limit(p.Limit).Order(p.Sort)
	       }).Find(&users)
	*/
	return users
}

func (dm *SqliteDataManager) GetAgents(offset int, limit int) []model.Agent {
	var agents []model.Agent
	dm.db.Offset(offset).Limit(limit).Find(&agents)
	return agents
}

func (dm *SqliteDataManager) GetAgentByParams(p *model.Params) []model.Agent {
	var agents []model.Agent
	error := dm.db.Select(p.Select).
		Where(p.FilterExp, p.FilterArgs).
		Offset(p.Offset).
		Limit(p.Limit).
		Order(p.Sort).
		Find(&agents).Error
	if error != nil {
		fmt.Printf("query error: %s \n", error)
	}
	return agents
}
