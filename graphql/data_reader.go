package graphql

import (
    "gorm.io/gorm"

	"github.com/chenfei531/query-api/model"
)

type DataReader interface {
	GetUserById(id int) model.User
	GetAgents(offset int, limit int) []model.Agent
}

type GraphDataReader struct {
	db *gorm.DB
}

func NewGraphDataReader(db *gorm.DB) *GraphDataReader {
	return &GraphDataReader{db}
}

func (dm *GraphDataReader) GetUserById(id int) model.User {
	var user model.User
	dm.db.Preload("Agents").First(&user, id)
	return user
}

func (dm *GraphDataReader) GetAgents(offset int, limit int) []model.Agent {
	var agents []model.Agent
	dm.db.Offset(offset).Limit(limit).Find(&agents)
	return agents
}
