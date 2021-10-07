package main

import (
	"fmt"
	"github.com/chenfei531/query-api/data"
	"github.com/chenfei531/query-api/query/graphql"
)

func main() {
    dm := data.NewSqliteDataManager()
	graphql.Init(dm)
	query := `
    {
        user(id: 2){
            id,
			name,
			agents(offset:2 limit:2){
				id,
				price,
				name
			},
        }
    }
    `
	result := graphql.Execute(query)
	fmt.Printf("%s \n", result)
}
