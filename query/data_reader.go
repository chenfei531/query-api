package query

import (
	"github.com/chenfei531/query-api/model"
)

type DataReader interface {
	GetData(root *Node) (string, error)
	//interface for graphql
	GetUserById(id int) model.User
	GetAgents(offset int, limit int) []model.Agent
}
