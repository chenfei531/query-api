package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chenfei531/query-api/model"
)

func HandleGraphQL(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("test\n")
	user := model.User{}
	if err := json.NewEncoder(w).Encode(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
}

/*
func getDBQuery(r *http.Request) (*rql.Params, error) {
	var (
		b   []byte
		err error
	)
	if v := r.URL.Query().Get(queryParam); v != "" {
		b, err = base64.StdEncoding.DecodeString(v)
	} else {
		b, err = ioutil.ReadAll(io.LimitReader(r.Body, 1<<12))
	}
	if err != nil {
		return nil, err
	}
	return queryParser.Parse(b)
}
*/

func main() {
	http.HandleFunc("/query", HandleGraphQL)
	http.ListenAndServe(":8080", nil)
}
