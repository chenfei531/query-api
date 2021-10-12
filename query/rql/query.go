package rql

import (
	"fmt"

	"github.com/chenfei531/query-api/data"
	"github.com/chenfei531/query-api/model"
)

/*
type ParamSchema interface {
    GetQueryParams(request string) (*model.Params, error)
}
*/

func GetQueryParamsByName(name string, request string) (*model.Params, error) {
	data, error := model.GetObjectByName(name)
	if nil != error {
		return nil, error
	}
	queryParser := MustNewParser(Config{Model: data, FieldSep: "."})
	p, error := queryParser.Parse([]byte(request))
	return (*model.Params)(p), error
}

func Do(dm data.DataManager, resource string, query_str string) string {
	p, error := GetQueryParamsByName(resource, query_str)
	if error != nil {
		fmt.Printf("%s \n", error)
	}
	r, error := dm.GetDataByNameAndParams(resource, p)
	if error != nil {
		fmt.Printf("%s \n", error)
	}
	return r
}
