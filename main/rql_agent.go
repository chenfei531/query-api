package main

import (
	"fmt"
	"github.com/chenfei531/query-api/data"
	"github.com/chenfei531/query-api/query/rql"
)

func main() {
	s := `
    {
      "limit": 3,
      "offset": 0,
      "filter": {
        "id": {"$gt": 9}
      },
      "sort": ["+name"],
      "select": ["name"]
    }
    `
	dm := data.NewSqliteDataManager()
	resp := rql.Do(dm, "Agent", s)
	fmt.Printf("%s \n", resp)
}
