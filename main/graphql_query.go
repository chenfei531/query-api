package main

import (
	"os"
	"fmt"
	"io/ioutil"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/chenfei531/query-api/graphql"
)

func main() {
	argsWithoutProg := os.Args[1:]
	var inputfile string
	if len(argsWithoutProg) == 0 {
		inputfile = "testquery/graphql_user"
	} else {
		inputfile = argsWithoutProg[0]
	}

	db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	dm := graphql.NewGraphDataReader(db)
	graphql.Init(dm)
	buf, _ := ioutil.ReadFile(inputfile)
	query := string(buf)
	result := graphql.Execute(query)
	fmt.Printf("%s \n", result)
}
