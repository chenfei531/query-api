package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/chenfei531/query-api/model"
)

type GormDataManager struct {
	db *gorm.DB
}

type Node struct {
	Name       string
	PrimaryKey string
	ContainPK  bool
	Params     model.Params
	Children   []*Node
}

func NewGormDataManager() *GormDataManager {
	db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	return &GormDataManager{db}
}

func getPrimayKey(data interface{}) string {
	return "ID"
}

func buildTreeNode(fields []string, index int, currentNode *Node, currentData interface{}) error {
	dataType := reflect.TypeOf(currentData)
	primaryKey := getPrimayKey(currentData)
	for i := index; i < len(fields); i++ {
		fieldName := fields[i]
		//handle nested fields

		//handle flat fields
		_, found := dataType.FieldByName(fieldName)
		if found {
			currentNode.Params.Select = append(currentNode.Params.Select, fieldName)
			if primaryKey == fieldName {
				currentNode.ContainPK = true
			}
		} else {
			_, found = dataType.FieldByName(fieldName + "s")
			if !found {
				return errors.New("field not found")
			}
			//validation
			node := Node{Name: fieldName, Params: model.Params{}, ContainPK: false}
			if nil == currentNode.Children {
				//necessary?
				currentNode.Children = make([]*Node, 0, 4)
			}
			currentNode.Children = append(currentNode.Children, &node)
		}
	}
	//add pk if not in select
	if !currentNode.ContainPK {
		currentNode.Params.Select = append(currentNode.Params.Select, primaryKey)
	}

	return nil
}

func buildQueryTree(name string, params *model.Params) (*Node, error) {
	rp := model.Params{Limit: params.Limit, Offset: params.Offset, Sort: params.Sort}
	rp.FilterExp = params.FilterExp
	rp.FilterArgs = params.FilterArgs
	//parse select into tree structure
	root := Node{Name: name, Params: rp, ContainPK: false}

	data, error := model.GetObjectByName(name)
	if nil != error {
		return nil, error
	}
	sorted_fields := params.Select
	sort.Strings(sorted_fields)
	fmt.Printf("sorted: %s\n", sorted_fields)
	error = buildTreeNode(sorted_fields, 0, &root, data)
	if nil != error {
		return nil, error
	}

	/*
		emptyParams := model.Params{}
		agentNode := Node{Name: "Agent", Params: emptyParams}
		//targetNode := Node{Name: "Target", Params: emptyParams}
		//agentNode.Children = []*Node{&targetNode}
		root.Children = []*Node{&agentNode}
	*/

	b, _ := json.MarshalIndent(root, "", "    ")
	fmt.Printf("tree:%s\n", b)
	return &root, nil
}

func (dm *GormDataManager) getDataByParams(data interface{}, p *model.Params) (interface{}, error) {
	dataType := reflect.TypeOf(data)
	result := reflect.New(reflect.SliceOf(dataType)).Interface()
	tx := dm.db.Model(data)
	if len(p.FilterExp) > 0 {
		//sometimes need to be p.FilterArgs..., why?
		tx.Where(p.FilterExp, p.FilterArgs...)
	}
	if 0 != p.Offset {
		tx.Offset(p.Offset)
	}
	if 0 != p.Limit {
		tx.Limit(p.Limit)
	}
	if len(p.Sort) > 0 {
		tx.Order(p.Sort)
	}
	if len(p.Select) > 0 {
		tx.Select(strings.Join(p.Select, ", "))
	}
	error := tx.Find(result).Error
	return result, error
}

func (dm *GormDataManager) getNestedData(node *Node) (interface{}, error) {
	name := node.Name
	params := node.Params
	data, error := model.GetObjectByName(name)
	if nil != error {
		return "", error
	}
	//get params
	//fmt.Printf("params:%s\n", params)
	results, error := dm.getDataByParams(data, &params)
	//fmt.Printf("results:%s\n", results)
	if nil != error {
		return "", error
	}

	var idList []uint64
	idIndexMap := make(map[uint64]uint)
	resultsReflect := reflect.ValueOf(results).Elem()
	for i := 0; i < resultsReflect.Len(); i++ {
		pkField := resultsReflect.Index(i).FieldByName("ID")
		id := pkField.Uint()
		idList = append(idList, id)
		idIndexMap[id] = uint(i)
		//remove primary key of parent if necessary
		//if !sp.containPK {
		//	pkField.SetUint(0)
		//}
	}

	for _, child := range node.Children {
		fkFieldName := name + "ID"
		fkDBName := name + "_" + "ID"
		childParamsPtr := &(child.Params)
		childParamsPtr.FilterExp = fkDBName + " IN (?)"
		fkArgs := make([]interface{}, 0, 8)
		for i := 0; i < len(idList); i++ {
			fkArgs = append(fkArgs, idList[uint(i)])
		}
		childParamsPtr.FilterArgs = make([]interface{}, 0, 2)
		childParamsPtr.FilterArgs = append(childParamsPtr.FilterArgs, fkArgs)

		children, _ := dm.getNestedData(child)
		childrenReflect := reflect.ValueOf(children).Elem()

		for i := 0; i < childrenReflect.Len(); i++ {
			childFieldName := child.Name + "s"
			foreignId := childrenReflect.Index(i).FieldByName(fkFieldName).Uint()
			parentIndex := int(idIndexMap[foreignId])
			parent := reflect.ValueOf(results).Elem().Index(parentIndex)
			childrenField := parent.FieldByName(childFieldName)
			buf := childrenField
			//remove unwanted fields
			childReflect := reflect.ValueOf(children).Elem().Index(i)
			childReflect.FieldByName(fkFieldName).SetUint(0)

			//append children to parents
			buf = reflect.Append(buf, childReflect)
			childrenField.Set(buf)
		}
	}
	return results, nil
}

func (dm *GormDataManager) GetDataByNameAndParams(name string, p *model.Params) (string, error) {
	root, error := buildQueryTree(name, p)
	if nil != error {
		return "", error
	}
	results, _ := dm.getNestedData(root)
	b, _ := json.MarshalIndent(results, "", "    ")
	return string(b), nil
}

func (dm *GormDataManager) GetUserById(id int) model.User {
	var user model.User
	dm.db.Preload("Agents").First(&user, id)
	return user
}

func (dm *GormDataManager) GetAgents(offset int, limit int) []model.Agent {
	var agents []model.Agent
	dm.db.Offset(offset).Limit(limit).Find(&agents)
	return agents
}
