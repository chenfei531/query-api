package query

import (
	"errors"
	"github.com/chenfei531/query-api/model"
	"github.com/chenfei531/query-api/query/rql"
)

type Node struct {
	Name       string
	PrimaryKey string
	ContainPK  bool
	Params     rql.Params //TODO: to ptr?
	Children   []*Node
}

func GetObjectByName(name string) (interface{}, error) {
	switch name {
	case "User":
		return model.User{}, nil
	case "Agent":
		return model.Agent{}, nil
	case "Target":
		return model.Target{}, nil
	case "MonitorLog":
		return model.MonitorLog{}, nil
	case "EventLog":
		return model.EventLog{}, nil
	default:
		return nil, errors.New("Type Not Found")
	}
}
