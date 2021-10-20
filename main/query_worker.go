package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/chenfei531/query-api/query"
)

func main() {
	buf, _ := ioutil.ReadFile("testdata/nested_query")
	array := strings.SplitAfterN(string(buf), "\n", 2)
	resource := strings.TrimSpace(array[0])
	queryStr := array[1]
	qw := query.NewQueryWorker()
	resp, error := qw.Query(resource, queryStr)
	if nil != error {
		fmt.Printf("%s \n", error)
		return
	}
	fmt.Printf("%s \n", resp)
}
