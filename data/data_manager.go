package data

import (
    "github.com/chenfei531/query-api/model"
)

type DataManager interface {
    GetUserById(id int) model.User
    GetUserByParams(p *model.Params) []model.User

    GetAgents(offset int, limit int) []model.Agent
    GetAgentByParams(p *model.Params) []model.Agent
}