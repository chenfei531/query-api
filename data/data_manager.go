package data

import (
	"github.com/chenfei531/query-api/model"
)

type DataManager interface {
	GetUserById(id int) model.User
	GetAgents(offset int, limit int) []model.Agent

	GetDataByNameAndParams(name string, p *model.Params) (string, error)

	getDataByParams(data interface{}, p *model.Params) (interface{}, error)
	getDataByNameAndParams(name string, p *model.Params) (interface{}, error)
}
