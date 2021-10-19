package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/chenfei531/query-api/data"
	"github.com/chenfei531/query-api/query/rql"
)

func main() {
	buf, _ := ioutil.ReadFile("test/data/nested_query")
	array := strings.SplitAfterN(string(buf), "\n", 2)
	resource := strings.TrimSpace(array[0])
	queryStr := array[1]
	dm := data.NewGormDataManager()
	resp := rql.Do(dm, resource, queryStr)
	//resp = ""
	fmt.Printf("%s \n", resp)
}
