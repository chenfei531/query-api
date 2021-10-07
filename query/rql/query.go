package rql

import (
	"encoding/json"
	"fmt"

	"github.com/chenfei531/query-api/data"
	"github.com/chenfei531/query-api/model"
)

func Do(dm *data.DataManager, query_str string) string {
	p, error := GetQueryParams(query_str)
	if error != nil {
		fmt.Printf("%s \n", error)
	}

	agents := dm.GetAgentByParams(p.Select, p.FilterExp, p.FilterArgs, p.Offset, p.Limit, p.Sort)
	b, err := json.MarshalIndent(agents, "", "    ")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return (string(b))
}

func GetQueryParams(request string) (*Params, error) {
	queryParser := MustNewParser(Config{
		Model:    model.Agent{}, // TODO: move to parameter
		FieldSep: ".",
	})
	return queryParser.Parse([]byte(request))
}
