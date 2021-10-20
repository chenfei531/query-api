package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/chenfei531/query-api/query"
)

func main() {
	argsWithoutProg := os.Args[1:]
	var inputfile string
	if len(argsWithoutProg) == 0 {
		inputfile = "testquery/nested_query"
	} else {
		inputfile = argsWithoutProg[0]
	}

	buf, _ := ioutil.ReadFile(inputfile)
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
