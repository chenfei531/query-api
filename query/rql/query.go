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
	var b []byte
	var err error
	switch(resource) {
	case "User":
		user := model.User{}
		p, error := GetQueryParams(&user, query_str)
		if error != nil {
			fmt.Printf("%s \n", error)
		}
		users := make([]model.User, 0)
		dm.GetDataByParams(&users, p)
		b, err = json.MarshalIndent(users, "", "    ")
	case "Agent":
		agent := model.Agent{}
		p, error := GetQueryParams(&agent, query_str)
		if error != nil {
			fmt.Printf("%s \n", error)
		}
		agents := make([]model.Agent, 0)
		dm.GetDataByParams(&agents, p)
		b, err = json.MarshalIndent(agents, "", "    ")
	default:
		return ""
	}
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return (string(b))
}
