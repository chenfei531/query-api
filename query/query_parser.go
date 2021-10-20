package query

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/chenfei531/query-api/query/rql"
)

type QueryParser struct {
}

func NewQueryParser() *QueryParser {
	return &QueryParser{}
}

//TODO: get pk/index by reflect tags
func (qp *QueryParser) getPrimayKey(data interface{}) string {
	return "ID"
}

func (qp *QueryParser) buildTreeNode(fields []string, index *int, currNode *Node, currPrefix string, currData interface{}) error {
	dataType := reflect.TypeOf(currData)
	primaryKey := qp.getPrimayKey(currData)
	for ; *index < len(fields); (*index)++ {
		i := *index
		endOfCurrNode := !strings.HasPrefix(fields[i], currPrefix)
		if endOfCurrNode {
			//add pk if not in select
			if !currNode.ContainPK {
				currNode.Params.Select = append(currNode.Params.Select, primaryKey)
			}
			(*index)--
			return nil
		}
		//strip prefix
		prefixLen := len(currPrefix)
		fieldName := fields[i][prefixLen:]

		elem := strings.SplitN(fieldName, ".", 2)
		if len(elem) > 1 {
			//handle nested field
			name := elem[0]
			prefix := currPrefix + name + "."
			var newNodePtr *Node
			l := len(currNode.Children)
			if l > 0 && currNode.Children[l-1].Name == name {
				newNodePtr = currNode.Children[l-1]
			} else {
				newNodePtr = &Node{Name: name, Params: rql.Params{}, ContainPK: false}
				currNode.Children = append(currNode.Children, newNodePtr)
			}

			data, error := GetObjectByName(name)
			if nil != error {
				return error
			}
			error = qp.buildTreeNode(fields, index, newNodePtr, prefix, data)
			if nil != error {
				return error
			}
		} else {
			//handle flat field
			_, found := dataType.FieldByName(fieldName)
			if found {
				currNode.Params.Select = append(currNode.Params.Select, fieldName)
				if primaryKey == fieldName {
					currNode.ContainPK = true
				}
			} else {
				_, found = dataType.FieldByName(fieldName + "s")
				if !found {
					return errors.New(fmt.Sprintf("field not found: %s", fieldName))
				}
				//TODO: validation
				node := Node{Name: fieldName, Params: rql.Params{}, ContainPK: true}
				node.Params.Select = append(node.Params.Select, "*")
				currNode.Children = append(currNode.Children, &node)
			}
		}
	}
	//add pk if not in select
	if !currNode.ContainPK {
		currNode.Params.Select = append(currNode.Params.Select, primaryKey)
	}

	return nil
}

func (qp *QueryParser) buildQueryTree(name string, params *rql.Params) (*Node, error) {
	rp := rql.Params{Limit: params.Limit, Offset: params.Offset, Sort: params.Sort}
	rp.FilterExp = params.FilterExp
	rp.FilterArgs = params.FilterArgs
	//parse select into tree structure
	root := Node{Name: name, Params: rp, ContainPK: false}

	data, error := GetObjectByName(name)
	if nil != error {
		return nil, error
	}
	sorted_fields := params.Select
	sort.Strings(sorted_fields)
	index := 0
	error = qp.buildTreeNode(sorted_fields, &index, &root, "", data)
	if nil != error {
		return nil, error
	}
	return &root, nil
}

func (qp *QueryParser) GetQueryTree(name string, request string) (*Node, error) {
	data, error := GetObjectByName(name)
	if nil != error {
		return nil, error
	}
	//TODO: cache Parser
	rqlParser := rql.MustNewParser(rql.Config{Model: data, FieldSep: "."})
	params, error := rqlParser.Parse([]byte(request))
	if nil != error {
		return nil, error
	}
	return qp.buildQueryTree(name, params)
}
