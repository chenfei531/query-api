package main

import (
    "fmt"
    "reflect"

    "github.com/chenfei531/query-api/model"
)

func Append(dest interface{}) {
    at := reflect.TypeOf(dest)
    t := at.Elem().Field(0).Tag
    fmt.Printf("%s\n", t)

    //reflect.MakeSlice(dest.type)
}

func main() {
    var users []model.User
    //users = append(users, model.User{})
    Append(users)
    fmt.Printf("%d\n", len(users))
}