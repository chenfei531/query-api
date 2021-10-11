package data

import (
	"github.com/chenfei531/query-api/model"
)

type DataManager interface {
	GetUserById(id int) model.User
	GetAgents(offset int, limit int) []model.Agent

	GetDataByParams(data interface{}, p *model.Params)
}
