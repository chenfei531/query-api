package data

import (
	"encoding/json"
	//"errors"
	"reflect"
	"strings"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/chenfei531/query-api/model"
)

type GormDataManager struct {
	db *gorm.DB
}

type Node struct {
    Name string
    Params model.Params
	Children []*Node
}

func NewGormDataManager() *GormDataManager {
	db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	return &GormDataManager{db}
}

func buildQueryTree(name string, params *model.Params) *Node {
    fmt.Printf("%s, %s\n", name, params)
	p := model.Params{}
	root := Node{Name: "User", Params: p}
	agentNode := Node{Name: "Agent", Params: p}
	targetNode := Node{Name: "Target", Params: p}
	agentNode.Children = []*Node{&targetNode}
	root.Children = []*Node{&agentNode}
	fmt.Printf("%s\n", agentNode)
    return &root
}

func (dm *GormDataManager) getDataByParams(data interface{}, p *model.Params) (interface{}, error) {
	dataType := reflect.TypeOf(data)
	result := reflect.New(reflect.SliceOf(dataType)).Interface()
	tx := dm.db.Model(data)
	if len(p.FilterExp) > 0 {
		tx.Where(p.FilterExp, p.FilterArgs)
	}
	if 0 != p.Offset {
		tx.Offset(p.Offset)
	}
	if 0 != p.Limit {
		tx.Limit(p.Limit)
	}
	if len(p.Sort) > 0 {
		tx.Order(p.Sort)
		fmt.Printf("%s\n", p.Sort)
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
	results, error := dm.getDataByParams(data, &params)
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
        children, _ := dm.getNestedData(child)
        childrenReflect := reflect.ValueOf(children).Elem()
		fkFieldName := name + "ID"
		//fkDBName := name + "_" + "ID"
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
    root := buildQueryTree(name, p)
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
