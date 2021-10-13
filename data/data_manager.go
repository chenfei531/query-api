package data

import (
	"github.com/chenfei531/query-api/model"
)

type DataManager interface {
	//interface for rql
	GetDataByNameAndParams(name string, p *model.Params) (string, error)
	getDataByParams(data interface{}, p *model.Params) (interface{}, error)
	//interface for graphql
	GetUserById(id int) model.User
	GetAgents(offset int, limit int) []model.Agent
}
