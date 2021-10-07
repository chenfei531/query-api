package main

import (
	"fmt"
	"io/ioutil"

	"github.com/chenfei531/query-api/data"
	"github.com/chenfei531/query-api/query/rql"
)

func main() {
	buf, _ := ioutil.ReadFile("test/data/rql_agent")
	query := string(buf)
	dm := data.NewSqliteDataManager()
	resp := rql.Do(dm, "Agent", query)
	fmt.Printf("%s \n", resp)
}
