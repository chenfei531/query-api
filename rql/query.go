package rql

import (
    "fmt"
    "encoding/json"

    "github.com/chenfei531/query-api/data"
)

func Do(query_str string) string {
    users := data.GetUserWithAgent()
    b, err := json.Marshal(users)
    if err != nil {
        fmt.Println(err)
        return ""
    }
    return (string(b))
}