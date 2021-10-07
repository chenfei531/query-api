package main

import (
	"fmt"
	"io/ioutil"

	"github.com/chenfei531/query-api/data"
	"github.com/chenfei531/query-api/query/graphql"
)

func main() {
	dm := data.NewSqliteDataManager()
	graphql.Init(dm)
	buf, _ := ioutil.ReadFile("test/data/graphql_agent")
	query := string(buf)
	result := graphql.Execute(query)
	fmt.Printf("%s \n", result)
}
