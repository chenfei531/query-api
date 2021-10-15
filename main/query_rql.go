package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/chenfei531/query-api/data"
	"github.com/chenfei531/query-api/query/rql"
)

func main() {
	buf, _ := ioutil.ReadFile("test/data/query_rql")
	array := strings.SplitAfterN(string(buf), "\n", 2)
	resource := strings.TrimSpace(array[0])
	queryStr := array[1]
	dm := data.NewSqliteDataManager()
	resp := rql.Do(dm, resource, queryStr)
	fmt.Printf("%s \n", resp)
}
