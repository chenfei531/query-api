package main

import (
	"fmt"
	"github.com/chenfei531/query-api/data"
	"github.com/chenfei531/query-api/query/graphql"
	"github.com/chenfei531/query-api/query/rql"
)

func main() {
	s := `
    {
      "limit": 10,
      "offset": 0,
      "filter": {
        "price": {"$gt": 900}
      },
      "sort": ["+name"],
      "select": ["name"]
    }
    `
	dm := data.NewDataManager()
	resp := rql.Do(dm, s)
	fmt.Printf("%s \n", resp)
	fmt.Printf("-------\n")
	graphql.Init(dm)
	query := `
    {
        agents(offset:2 limit:2){
            id
        }
    }
    `
	result := graphql.Execute(query)
	fmt.Printf("%s \n", result)
}
