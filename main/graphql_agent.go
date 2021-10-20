package main

import (
	"fmt"
	"io/ioutil"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/chenfei531/query-api/query"
	"github.com/chenfei531/query-api/query/graphql"
)

func main() {
	db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	dm := query.NewGormDataReader(db)
	graphql.Init(dm)
	buf, _ := ioutil.ReadFile("testdata/graphql_agent")
	query := string(buf)
	result := graphql.Execute(query)
	fmt.Printf("%s \n", result)
}
