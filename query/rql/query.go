package rql

import (
	"encoding/json"
	"fmt"

	"github.com/chenfei531/query-api/data"
	"github.com/chenfei531/query-api/model"
)

/*
type ParamSchema interface {
    GetQueryParams(request string) (*model.Params, error)
}
*/

func GetQueryParams(emptyObj interface{}, request string) (*model.Params, error) {
	queryParser := MustNewParser(Config{
		// TODO: use reflect
		Model:    emptyObj,
		FieldSep: ".",
	})
	p, error := queryParser.Parse([]byte(request))
	return (*model.Params)(p), error
}

func Do(dm data.DataManager, resource string, query_str string) string {
	//TODO: use reflect instead of 'if'
	if resource == "Agent" {
		p, error := GetQueryParams(model.Agent{}, query_str)
		if error != nil {
			fmt.Printf("%s \n", error)
		}

		agents := dm.GetAgentByParams(p)
		b, err := json.MarshalIndent(agents, "", "    ")
		if err != nil {
			fmt.Println(err)
			return ""
		}
		return (string(b))
	}
	if resource == "User" {
		p, error := GetQueryParams(model.User{}, query_str)
		if error != nil {
			fmt.Printf("%s \n", error)
		}

		users := dm.GetUserByParams(p)
		b, err := json.MarshalIndent(users, "", "    ")
		if err != nil {
			fmt.Println(err)
			return ""
		}
		return (string(b))
	}
	return ""
}
