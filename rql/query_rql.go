package rql

import (
    "github.com/chenfei531/query-api/model"
)

func GetQueryParams(request string) (*Params, error) {
    queryParser := MustNewParser(Config{
		Model:    model.Agent{},
		FieldSep: ".",
	})
    return queryParser.Parse([]byte(request))
}
