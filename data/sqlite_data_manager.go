package data

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"reflect"

	"github.com/chenfei531/query-api/model"
)

type SqliteDataManager struct {
	db *gorm.DB
}

type selectParams struct {
	 fields []string
	 nestedField map[string][]string
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

func (sp selectParams) getSelectParams(selectParams []string) {

}

func (dm *SqliteDataManager) GetDataByParams(data interface{}, p *model.Params) {
	t := reflect.TypeOf(data).Elem().Elem().Field(0).Tag
	fmt.Printf("%s\n", t)

	/*
	Select: must contain ID, otherwise preload will not work
	select does not work with preloaded fields: https://github.com/go-gorm/gorm/issues/4015
	*/
	tx := dm.db.Select(p.Select).
		Where(p.FilterExp, p.FilterArgs).
		Offset(p.Offset).
		Limit(p.Limit).
		Order(p.Sort)
	tx = tx.Preload("Agents")
	tx.Find(data)
}

func (dm *SqliteDataManager) GetAgents(offset int, limit int) []model.Agent {
	var agents []model.Agent
	dm.db.Offset(offset).Limit(limit).Find(&agents)
	return agents
}
