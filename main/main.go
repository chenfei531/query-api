package main

import (
    "fmt"
    //"github.com/chenfei531/query-api/data"
    "github.com/chenfei531/query-api/rql"
)

func main() {
    //users := data.GetUserWithAgent()
    //fmt.Printf("%s \n", users[1].Agents[1].Name)
    s := `
    {
      "limit": 25,
      "offset": 0,
      "filter": {
        "price": {"$gt": 100}
      },
      "sort": ["+name"]
    }
    `
    /*
    r, error := rql.GetQueryParams(s)
    if error != nil {
        fmt.Printf("%s \n", error)
    }
    fmt.Printf("%s \n", r)
    */
    resp := rql.Do(s)
    fmt.Printf("%s \n", resp)
}